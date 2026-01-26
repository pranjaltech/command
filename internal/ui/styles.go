package ui

import "github.com/charmbracelet/lipgloss"

// Color palette for consistent styling across the UI
var (
	// SuccessStyle is used for checkmarks and success messages
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")) // green

	// ErrorStyle is used for error messages and X marks
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red

	// MutedStyle is used for hints and secondary text
	MutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // gray

	// BoldStyle is used for emphasis
	BoldStyle = lipgloss.NewStyle().Bold(true)

	// TitleStyle is used for section titles
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99")) // purple

	// BoxStyle is used for framed content boxes
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	// SelectedStyle is used for currently selected items
	SelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")) // pink/magenta
)

// Checkmark returns a green checkmark symbol
func Checkmark() string {
	return SuccessStyle.Render("\u2713")
}

// Cross returns a red X symbol
func Cross() string {
	return ErrorStyle.Render("\u2717")
}

// MaskAPIKey masks an API key, showing only the last 4 characters
func MaskAPIKey(key string) string {
	if key == "" {
		return MutedStyle.Render("(not set)")
	}
	if len(key) <= 4 {
		return MutedStyle.Render("\u2022\u2022\u2022\u2022")
	}
	masked := ""
	for i := 0; i < len(key)-4; i++ {
		masked += "\u2022"
	}
	return masked + key[len(key)-4:]
}
