package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"ghprofile/github"
	"ghprofile/ui"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`
	ghprofile: Pretty GitHub profile viewer

	Usage:
		ghprofile [flags]

	Flags:
		-u, --user        GitHub username to fetch (default: dayvster)
		-n                How many top repos to show (default: 5)
		--no-icons        Disable icons in the output
		--no-border       Remove card border from output
		--no-style        Remove all styles from output
		--no-demo         Do not fall back to demo data on fetch error; exit instead
		--demo            Force demo data (skip network and cache)
		-h, --help        Show this help message
`)
	}
	userLong := flag.String("user", "dayvster", "GitHub username to fetch")
	userShort := flag.String("u", "", "GitHub username (shorthand)")
	topN := flag.Int("n", 5, "How many top repos to show")
	noIcons := flag.Bool("no-icons", false, "Disable icons in the output")
	noDemo := flag.Bool("no-demo", false, "Do not fall back to demo data on fetch error; exit instead")
	demo := flag.Bool("demo", false, "Force demo data (skip network and cache)")
	noBorder := flag.Bool("no-border", false, "Remove card border from output")
	noStyle := flag.Bool("no-style", false, "Remove all styles from output")
	flag.Parse()

	user := *userLong
	if *userShort != "" {
		user = *userShort
	}

	gh := &github.Github{Client: http.DefaultClient}

	if *demo {
		p, repos := github.DemoProfile(github.DemoProfileConfig{Username: user})
		showIcons := !*noIcons
		ui.PrintProfile(p, repos, *topN, showIcons, *noBorder, *noStyle)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	p, repos, err := gh.FetchProfileWithRepos(ctx, user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: fetch failed: %v\n", err)
		if strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "API rate limit") {
		}
		if *noDemo {
			if cp, cr, cerr := github.TryLoadCache(user); cerr == nil {
				fmt.Fprintf(os.Stderr, "loaded cached profile for %s\n", user)
				p, repos = cp, cr
			} else {
				fmt.Fprintf(os.Stderr, "no cache available and --no-demo set; exiting\n")
				os.Exit(1)
			}
		} else {
			if cp, cr, cerr := github.TryLoadCache(user); cerr == nil {
				fmt.Fprintf(os.Stderr, "warning: fetch failed — using cached data for %s\n", user)
				p, repos = cp, cr
			} else {
				fmt.Fprintf(os.Stderr, "warning: fetch failed (%v) — falling back to demo data for %s\n", err, user)
				p, repos = github.DemoProfile(github.DemoProfileConfig{Username: user})
			}
		}
	} else {
		if err := github.SaveCache(user, p, repos); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to save cache: %v\n", err)
		}
	}

	showIcons := !*noIcons
	ui.PrintProfile(p, repos, *topN, showIcons, *noBorder, *noStyle)
}
