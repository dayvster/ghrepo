package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

const DefaultUsername = "dayvster"

type rewriteTransport struct {
	base   http.RoundTripper
	target *url.URL
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.base == nil {
		t.base = http.DefaultTransport
	}
	r2 := req.Clone(req.Context())
	r2.URL.Scheme = t.target.Scheme
	r2.URL.Host = t.target.Host
	r2.Host = t.target.Host
	return t.base.RoundTrip(r2)
}

func TestFetchProfileAndRepos(t *testing.T) {
	if os.Getenv("REAL_GITHUB") == "1" {
		t.Logf("running integration fetch for %s", DefaultUsername)
		ctx := context.Background()
		gh := &Github{}
		p, repos, err := gh.FetchProfileWithRepos(ctx, DefaultUsername)
		if err != nil {
			t.Fatalf("integration fetch error: %v", err)
		}
		totalStars := 0
		if p.TotalStars != nil {
			totalStars = *p.TotalStars
		}
		totalForks := 0
		if p.TotalForks != nil {
			totalForks = *p.TotalForks
		}
		avg := float32(0)
		if p.AvgStarsPerRepo != nil {
			avg = *p.AvgStarsPerRepo
		}
		t.Logf("Profile: %s (%s)", p.FullName, p.URL)
		t.Logf("Public repos: %d, total stars: %d, total forks: %d, avg stars/repo: %.2f", p.PublicReposAmount, totalStars, totalForks, avg)
		t.Logf("Fetched %d repos", len(repos))
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users/"+DefaultUsername, func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"login":            DefaultUsername,
			"avatar_url":       "http://avatar",
			"html_url":         "http://html",
			"name":             "Test User",
			"company":          "Acme",
			"blog":             "https://blog",
			"bio":              "bio",
			"twitter_username": "twt",
			"followers":        10,
			"following":        2,
			"created_at":       "2020-01-01T00:00:00Z",
			"hireable":         true,
			"email":            "t@example.com",
			"public_repos":     2,
			"public_gists":     1,
		}
		b, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	mux.HandleFunc("/users/"+DefaultUsername+"/repos", func(w http.ResponseWriter, r *http.Request) {
		repos := []map[string]interface{}{
			{"id": 1, "name": "r1", "full_name": DefaultUsername + "/r1", "html_url": "u1", "stargazers_count": 5, "forks_count": 1},
			{"id": 2, "name": "r2", "full_name": DefaultUsername + "/r2", "html_url": "u2", "stargazers_count": 3, "forks_count": 2},
		}
		b, _ := json.Marshal(repos)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	client := &http.Client{Transport: &rewriteTransport{target: u}}

	ctx := context.Background()
	gh := &Github{Client: client}
	runProfileAndReposTests(t, gh, ctx, DefaultUsername)

}

func runProfileAndReposTests(t *testing.T, gh *Github, ctx context.Context, username string) {
	p, err := gh.GetProfile(ctx, username)
	if err != nil {
		t.Fatalf("GetProfile error: %v", err)
	}
	if p.Name != username {
		t.Fatalf("expected Name %s got %q", username, p.Name)
	}

	repos, err := gh.GetRepos(ctx, username)
	if err != nil {
		t.Fatalf("GetRepos error: %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos got %d", len(repos))
	}

	p2, repos2, err := gh.FetchProfileWithRepos(ctx, username)
	if err != nil {
		t.Fatalf("FetchProfileWithRepos error: %v", err)
	}
	if p2.TotalStars == nil || *p2.TotalStars != 8 {
		t.Fatalf("expected total stars 8 got %v", p2.TotalStars)
	}
	if p2.TotalForks == nil || *p2.TotalForks != 3 {
		t.Fatalf("expected total forks 3 got %v", p2.TotalForks)
	}
	if p2.AvgStarsPerRepo == nil || *p2.AvgStarsPerRepo != 4.0 {
		t.Fatalf("expected avg 4.0 got %v", p2.AvgStarsPerRepo)
	}
	if len(repos2) != 2 {
		t.Fatalf("expected repos2 len 2 got %d", len(repos2))
	}
}

func TestGetReposPagination(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/"+DefaultUsername+"/repos", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		if page == "" || page == "1" {
			repos := make([]map[string]interface{}, 0, 100)
			for i := 1; i <= 100; i++ {
				repos = append(repos, map[string]interface{}{"id": i, "name": fmt.Sprintf("r%d", i), "full_name": DefaultUsername + fmt.Sprintf("/r%d", i), "stargazers_count": i, "forks_count": i % 3})
			}
			b, _ := json.Marshal(repos)
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		if page == "2" {
			repos := []map[string]interface{}{{"id": 101, "name": "r101", "full_name": DefaultUsername + "/r101", "stargazers_count": 101, "forks_count": 1}}
			b, _ := json.Marshal(repos)
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	client := &http.Client{Transport: &rewriteTransport{target: u}}

	gh := &Github{Client: client}
	runReposPaginationTest(t, gh, DefaultUsername)
}

func runReposPaginationTest(t *testing.T, gh *Github, username string) {
	repos, err := gh.GetRepos(context.Background(), username)
	if err != nil {
		t.Fatalf("GetRepos paginated error: %v", err)
	}
	if len(repos) != 101 {
		t.Fatalf("expected 101 repos from paginated handler, got %d", len(repos))
	}
}

func TestGetProfileNonOK(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/"+DefaultUsername, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	client := &http.Client{Transport: &rewriteTransport{target: u}}

	gh := &Github{Client: client}
	runProfileNonOKTest(t, gh, DefaultUsername)
}

func runProfileNonOKTest(t *testing.T, gh *Github, username string) {
	_, err := gh.GetProfile(context.Background(), username)
	if err == nil {
		t.Fatalf("expected error for non-OK response")
	}
	if !strings.Contains(err.Error(), "returned 404") {
		t.Fatalf("expected error mentioning returned 404, got: %v", err)
	}
}

func TestUserAgentHeader(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/users/"+DefaultUsername, func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		if ua != "ghprofile-client" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad ua"))
			return
		}
		res := map[string]interface{}{"login": DefaultUsername}
		b, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	u, _ := url.Parse(srv.URL)
	client := &http.Client{Transport: &rewriteTransport{target: u}}

	gh := &Github{Client: client}
	runUserAgentHeaderTest(t, gh, DefaultUsername)
}

func runUserAgentHeaderTest(t *testing.T, gh *Github, username string) {
	_, err := gh.GetProfile(context.Background(), username)
	if err != nil {
		t.Fatalf("GetProfile failed when testing User-Agent: %v", err)
	}
}
