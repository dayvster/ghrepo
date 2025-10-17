package ui

import (
	"context"
	"fmt"
	"time"

	"ghprofile/github"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	username string
	loading  bool
	spinner  spinner.Model
	profile  *github.Profile
	repos    []github.Repo
	err      error
	client   *github.Github
}

func New(username string, client *github.Github) tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return &model{
		username: username,
		loading:  true,
		spinner:  s,
		client:   client,
	}
}

type fetchMsg struct{}

type fetchResult struct {
	profile *github.Profile
	repos   []github.Repo
	err     error
}

func fetchCmd(username string, github *github.Github) tea.Cmd {
	return func() tea.Msg {
		p, repos, err := github.FetchProfileWithRepos(context.Background(), username)
		return fetchResult{profile: p, repos: repos, err: err}
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		time.Sleep(10 * time.Millisecond)
		return fetchMsg{}
	})
}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case fetchMsg:
		return m, tea.Batch(fetchCmd(m.username, m.client), m.spinner.Tick)
	case fetchResult:
		m.loading = false
		m.profile = msg.profile
		m.repos = msg.repos
		m.err = msg.err
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if m.err != nil {
		return Panel.Render(fmt.Sprintf("Error: %v", m.err))
	}
	if m.loading && m.profile == nil {
		return Panel.Render(fmt.Sprintf("%s Loading %s", TitleStyle.Render("ghprofile"), m.spinner.View()))
	}
	if m.profile == nil {
		return Panel.Render("No profile")
	}
	totalStars := 0
	totalForks := 0
	if m.profile.TotalStars != nil {
		totalStars = *m.profile.TotalStars
	}
	if m.profile.TotalForks != nil {
		totalForks = *m.profile.TotalForks
	}
	avg := 0.0
	if m.profile.AvgStarsPerRepo != nil {
		avg = float64(*m.profile.AvgStarsPerRepo)
	}

	out := fmt.Sprintf("%s\n\n", TitleStyle.Render(m.profile.FullName))
	out += fmt.Sprintf("%s %s\n", Subtle.Render("URL:"), m.profile.URL)
	out += fmt.Sprintf("%s %d\n", StatStyle.Render("Public repos:"), m.profile.PublicReposAmount)
	out += fmt.Sprintf("%s %d\n", StatStyle.Render("Total stars:"), totalStars)
	out += fmt.Sprintf("%s %d\n", StatStyle.Render("Total forks:"), totalForks)
	out += fmt.Sprintf("%s %.2f\n", StatStyle.Render("Avg stars/repo:"), avg)

	return Panel.Render(out)
}
