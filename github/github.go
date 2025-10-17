package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Github struct {
	Client *http.Client
}

type Profile struct {
	Name              string    `json:"name,omitempty"`
	AvatarURL         string    `json:"avatar_url,omitempty"`
	URL               string    `json:"url,omitempty"`
	FullName          string    `json:"full_name,omitempty"`
	Company           string    `json:"company,omitempty"`
	Blog              string    `json:"blog,omitempty"`
	Bio               string    `json:"bio,omitempty"`
	Twitter           string    `json:"twitter,omitempty"`
	FollowersAmount   int       `json:"followers_amount,omitempty"`
	FollowingAmount   int       `json:"following_amount,omitempty"`
	MemberSince       string    `json:"member_since,omitempty"`
	Hireable          bool      `json:"hireable,omitempty"`
	Email             string    `json:"email,omitempty"`
	PublicReposAmount int       `json:"public_repos_amount,omitempty"`
	PublicGistsAmount int       `json:"public_gists_amount,omitempty"`
	TotalStars        *int      `json:"total_stars,omitempty"`
	TotalForks        *int      `json:"total_forks,omitempty"`
	AvgStarsPerRepo   *float32  `json:"avg_stars_per_repo,omitempty"`
	Repos             []Repo    `json:"repos,omitempty"`
	Followers         []Profile `json:"followers,omitempty"`
}

type Repo struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	HTMLURL         string `json:"html_url,omitempty"`
	Description     string `json:"description,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	ForksCount      int    `json:"forks_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	Language        string `json:"language,omitempty"`
	Size            int    `json:"size,omitempty"`
	OpenIssuesCount int    `json:"open_issues_count,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	Fork            bool   `json:"fork,omitempty"`
	DefaultBranch   string `json:"default_branch,omitempty"`
}

func (gh *Github) doRequest(ctx context.Context, method, urlStr string) ([]byte, error) {
	client := gh.Client
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, method, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	req.Header.Set("User-Agent", "ghprofile-client")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: %s %s returned %d: %s", method, urlStr, resp.StatusCode, string(body))
	}
	return body, nil
}

func (gh *Github) GetProfile(ctx context.Context, username string) (*Profile, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	u := fmt.Sprintf("https://api.github.com/users/%s", url.PathEscape(username))
	body, err := gh.doRequest(ctx, http.MethodGet, u)
	if err != nil {
		return nil, err
	}
	var g struct {
		Login       string `json:"login"`
		AvatarURL   string `json:"avatar_url"`
		HTMLURL     string `json:"html_url"`
		Name        string `json:"name"`
		Company     string `json:"company"`
		Blog        string `json:"blog"`
		Bio         string `json:"bio"`
		Twitter     string `json:"twitter_username"`
		Followers   int    `json:"followers"`
		Following   int    `json:"following"`
		CreatedAt   string `json:"created_at"`
		Hireable    bool   `json:"hireable"`
		Email       string `json:"email"`
		PublicRepos int    `json:"public_repos"`
		PublicGists int    `json:"public_gists"`
	}
	if err := json.Unmarshal(body, &g); err != nil {
		return nil, fmt.Errorf("unmarshal profile: %w", err)
	}
	p := &Profile{
		Name:              g.Login,
		AvatarURL:         g.AvatarURL,
		URL:               g.HTMLURL,
		FullName:          g.Name,
		Company:           g.Company,
		Blog:              g.Blog,
		Bio:               g.Bio,
		Twitter:           g.Twitter,
		FollowersAmount:   g.Followers,
		FollowingAmount:   g.Following,
		MemberSince:       g.CreatedAt,
		Hireable:          g.Hireable,
		Email:             g.Email,
		PublicReposAmount: g.PublicRepos,
		PublicGistsAmount: g.PublicGists,
	}
	return p, nil
}

func (gh *Github) paginateRepos(ctx context.Context, username string) ([]Repo, error) {
	var all []Repo
	for page := 1; ; page++ {
		u := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", url.PathEscape(username), page)
		body, err := gh.doRequest(ctx, http.MethodGet, u)
		if err != nil {
			return nil, err
		}
		var repos []Repo
		if err := json.Unmarshal(body, &repos); err != nil {
			return nil, fmt.Errorf("unmarshal repos: %w", err)
		}
		if len(repos) == 0 {
			break
		}
		all = append(all, repos...)
		if len(repos) < 100 {
			break
		}
	}
	return all, nil
}

func (gh *Github) GetRepos(ctx context.Context, username string) ([]Repo, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	return gh.paginateRepos(ctx, username)
}

func (gh *Github) calcRepoStats(p *Profile, repos []Repo) {
	totalStars := 0
	totalForks := 0
	for _, r := range repos {
		totalStars += r.StargazersCount
		totalForks += r.ForksCount
	}
	if len(repos) > 0 {
		avg := float32(totalStars) / float32(len(repos))
		p.AvgStarsPerRepo = new(float32)
		*p.AvgStarsPerRepo = avg
	}
	p.TotalStars = new(int)
	*p.TotalStars = totalStars
	p.TotalForks = new(int)
	*p.TotalForks = totalForks
}

func (gh *Github) FetchProfileWithRepos(ctx context.Context, username string) (*Profile, []Repo, error) {
	p, err := gh.GetProfile(ctx, username)
	if err != nil {
		return nil, nil, err
	}
	repos, err := gh.GetRepos(ctx, username)
	if err != nil {
		return p, nil, err
	}
	gh.calcRepoStats(p, repos)
	return p, repos, nil
}
