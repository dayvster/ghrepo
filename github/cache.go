/**
This package provides caching functionality for GitHub user profiles and repositories.
It allows saving and loading profile data to and from a cache directory, improving
performance by avoiding repeated API calls to GitHub.

*/

package github

import (
	"encoding/json"
	"os"
	"time"
)

const (
	DefaultCacheDir = ".cache/ghprofile"
	CacheFileSuffix = ".json"
	CacheExpire     = 30 * time.Minute
)

type cacheEntry struct {
	Profile *Profile `json:"profile"`
	Repos   []Repo   `json:"repos"`
}

func CachePath(user string) (string, error) {
	dir := os.Getenv("XDG_CACHE_HOME")
	if dir == "" {
		dir = os.Getenv("HOME") + "/.cache"
	}
	base := dir + "/ghprofile"
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	return base + "/" + user + CacheFileSuffix, nil
}

func SaveCache(user string, p *Profile, repos []Repo) error {
	path, err := CachePath(user)
	if err != nil {
		return err
	}
	ce := cacheEntry{Profile: p, Repos: repos}
	b, err := json.MarshalIndent(ce, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func TryLoadCache(user string) (*Profile, []Repo, error) {
	path, err := CachePath(user)
	if err != nil {
		return nil, nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var ce cacheEntry
	if err := json.Unmarshal(b, &ce); err != nil {
		return nil, nil, err
	}
	return ce.Profile, ce.Repos, nil
}
