package ui

import "github.com/charmbracelet/lipgloss"

var (
	panelBg      = lipgloss.Color("#0f1724")
	cardBg       = lipgloss.Color("#141826")
	headerBg     = lipgloss.Color("#1f2335")
	brightFg     = lipgloss.Color("#c0caf5")
	accentCyan   = lipgloss.Color("#7dcfff")
	accentBlue   = lipgloss.Color("#7aa2f7")
	accentPurple = lipgloss.Color("#bb9af7")

	Panel = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).Width(88).
		BorderForeground(accentBlue).Foreground(brightFg)

	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(accentCyan)

	URLStyle = lipgloss.NewStyle().Foreground(accentBlue).Underline(true)

	Subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9aa5ff"))

	StatStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9ece6a"))

	Accent = lipgloss.NewStyle().Bold(true).Foreground(accentPurple)

	IconStyle = lipgloss.NewStyle().Bold(true).Foreground(accentCyan)

	RepoTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffb86b"))

	ValueStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9ece6a"))

	StatBox = lipgloss.NewStyle().Padding(0, 1).Bold(true).Foreground(lipgloss.Color("#9ece6a")).MarginRight(1)

	Badge = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("#ff9e64")).Bold(true).MarginRight(1)

	Divider = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7089"))
)
