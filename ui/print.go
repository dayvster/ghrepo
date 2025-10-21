package ui

import (
	"fmt"
	"sort"
	"strings"

	"ghprofile/github"

	"github.com/charmbracelet/lipgloss"
)

func PrintProfile(p *github.Profile, repos []github.Repo, topN int, showIcons bool, noBorder bool, noStyle bool, size string) {
	if p == nil {
		fmt.Println("No profile")
		return
	}

	var b strings.Builder
	title := TitleStyle.Render(fmt.Sprintf("%s", p.FullName))
	url := URLStyle.Render(p.URL)
	iconRender := func(s string) string {
		if !showIcons || s == "" {
			return ""
		}
		return IconStyle.Render(s)
	}
	if iconRender(IconUser) != "" {
		b.WriteString(iconRender(IconUser) + "  " + title + "  " + url + "\n")
	} else {
		b.WriteString(title + "  " + url + "\n")
	}
	if p.Bio != "" {
		b.WriteString(Subtle.Render(p.Bio) + "\n")
	}
	b.WriteString("\n")
	type statEntry struct{ icon, label, value string }
	totalStars := 0
	totalForks := 0
	if p.TotalStars != nil {
		totalStars = *p.TotalStars
	}
	if p.TotalForks != nil {
		totalForks = *p.TotalForks
	}
	avg := float32(0)
	if p.AvgStarsPerRepo != nil {
		avg = *p.AvgStarsPerRepo
	}

	stats := []statEntry{
		{iconRender(IconFollowers), "Followers:", fmt.Sprintf("%d", p.FollowersAmount)},
		{iconRender(IconFollowing), "Following:", fmt.Sprintf("%d", p.FollowingAmount)},
		{iconRender(IconRepo), "Public repos:", fmt.Sprintf("%d", p.PublicReposAmount)},
		{iconRender(IconGist), "Public gists:", fmt.Sprintf("%d", p.PublicGistsAmount)},
		{iconRender(IconStar), "Total stars:", fmt.Sprintf("%d", totalStars)},
		{iconRender(IconFork), "Total forks:", fmt.Sprintf("%d", totalForks)},
		{iconRender(IconStar), "Avg stars/repo:", fmt.Sprintf("%.2f", avg)},
	}
	maxLabel := 0
	for _, s := range stats {
		if len(s.label) > maxLabel {
			maxLabel = len(s.label)
		}
	}
	for _, s := range stats {
		if s.icon != "" {
			b.WriteString(s.icon + "  ")
		} else {
			b.WriteString("   ")
		}
		padded := fmt.Sprintf("%-*s", maxLabel, s.label)
		b.WriteString(Accent.Render(padded) + " " + ValueStyle.Render(s.value) + "\n")
	}

	langCount := map[string]int{}
	for _, r := range repos {
		lang := r.Language
		if lang == "" {
			continue // ignore repos with no language
		}
		langCount[lang]++
	}
	if len(langCount) > 0 {
		b.WriteString("\n")
		b.WriteString(Subtle.Render("Languages:") + "\n")
		type kv struct {
			k string
			v int
		}
		var kvs []kv
		for k, v := range langCount {
			kvs = append(kvs, kv{k, v})
		}
		sort.Slice(kvs, func(i, j int) bool { return kvs[i].v > kvs[j].v })
		for _, x := range kvs {
			icon := GetLangIcon(x.k)
			out := fmt.Sprintf("%s  %s: %d\n", iconRender(icon), Accent.Render(x.k), x.v)
			if icon == "" {
				out = fmt.Sprintf("   %s: %d\n", Accent.Render(x.k), x.v)
			}
			b.WriteString(out)
		}
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].StargazersCount > repos[j].StargazersCount })
	if topN > len(repos) {
		topN = len(repos)
	}
	if topN > 0 {
		b.WriteString("\n")
		if noStyle {
			b.WriteString("Top repos:\n")
		} else {
			b.WriteString(Subtle.Render("Top repos:") + "\n")
		}
		for i := 0; i < topN; i++ {
			r := repos[i]
			langIcon := GetLangIcon(r.Language)
			if noStyle {
				b.WriteString(fmt.Sprintf("%d. %s %s ★ %d  %s %d\n", i+1, r.FullName, r.Language, r.StargazersCount, IconFork, r.ForksCount))
				b.WriteString("  " + r.HTMLURL + "\n")
			} else {
				b.WriteString(fmt.Sprintf("%d. %s %s %s %d  %s %d\n", i+1, RepoTitle.Render(r.FullName), iconRender(langIcon), iconRender("★"), r.StargazersCount, iconRender(IconFork), r.ForksCount))
				b.WriteString("  " + URLStyle.Render(r.HTMLURL) + "\n")
			}
		}
	}

	out := b.String()
	if noStyle || noBorder {
		fmt.Print(out)
		return
	}

	// Determine width from size flag. 0 means no width constraint (full).
	width := 0
	switch strings.ToLower(size) {
	case "small":
		width = 48
	case "medium":
		width = 88
	case "large":
		width = 120
	case "full":
		width = 0
	default:
		width = 88
	}

	panel := Panel(width).Render(out)
	fmt.Println(lipgloss.NewStyle().Margin(1, 2).Render(panel))
}
