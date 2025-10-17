package github

import "fmt"

type DemoProfileConfig struct {
	Username string
}

func DemoProfile(cfg DemoProfileConfig) (*Profile, []Repo) {
	uname := cfg.Username
	if uname == "" {
		uname = "demo"
	}
	p := &Profile{
		FullName:          uname,
		URL:               fmt.Sprintf("https://github.com/%s", uname),
		Bio:               `This is demo data used when the GitHub API is unavailable.\n\nNote: The GitHub API is unavailable in this environment.`,
		FollowersAmount:   42,
		FollowingAmount:   7,
		PublicReposAmount: 5,
		PublicGistsAmount: 1,
	}
	repos := []Repo{
		{FullName: fmt.Sprintf("%s/repo-one", uname), HTMLURL: fmt.Sprintf("https://github.com/%s/repo-one", uname), StargazersCount: 10, ForksCount: 2, Language: "Go"},
		{FullName: fmt.Sprintf("%s/repo-two", uname), HTMLURL: fmt.Sprintf("https://github.com/%s/repo-two", uname), StargazersCount: 5, ForksCount: 1, Language: "TypeScript"},
		{FullName: fmt.Sprintf("%s/repo-three", uname), HTMLURL: fmt.Sprintf("https://github.com/%s/repo-three", uname), StargazersCount: 2, ForksCount: 0, Language: "Rust"},
	}
	totalStars := 0
	totalForks := 0
	for _, r := range repos {
		totalStars += r.StargazersCount
		totalForks += r.ForksCount
	}
	p.TotalStars = new(int)
	*p.TotalStars = totalStars
	p.TotalForks = new(int)
	*p.TotalForks = totalForks
	if len(repos) > 0 {
		tmp := float32(totalStars) / float32(len(repos))
		p.AvgStarsPerRepo = new(float32)
		*p.AvgStarsPerRepo = tmp
	}
	return p, repos
}
