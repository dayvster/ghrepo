package ui

import "strings"

var LangIcons = map[string]string{
	"typescript": "",
	"javascript": "",
	"go":         "",
	"python":     "",
	"rust":       "",
	"c++":        "󰙲",
	"c#":         "",
	"php":        "",
	"java":       "",
	"zig":        "",
	"shell":      "",
	"vue":        "",
	"odin":       "◎",
	"react":      "",
	"ruby":       "",
	"swift":      "",
	"kotlin":     "",
	"dart":       "",
	"elixir":     "",
	"haskell":    "",
}

func GetLangIcon(lang string) string {
	if lang == "" {
		return "(unknown)"
	}
	k := strings.ToLower(lang)
	if v, ok := LangIcons[k]; ok {
		return v
	}
	return "λ"
}
