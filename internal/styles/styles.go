package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Primary   = lipgloss.Color("63")  // Purple
	Secondary = lipgloss.Color("36")  // Cyan
	Success   = lipgloss.Color("42")  // Green
	Warning   = lipgloss.Color("226") // Yellow
	Danger    = lipgloss.Color("196") // Red
	Muted     = lipgloss.Color("240") // Gray
	Text      = lipgloss.Color("255") // White
	
	// Common Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true).
			Padding(1, 2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	// Panel Styles
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	ActivPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(1, 2)

	// List Styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true).
				PaddingLeft(2)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	// Method Styles
	MethodGET = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true).
			Width(7)

	MethodPOST = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")). // Blue
			Bold(true).
			Width(7)

	MethodPUT = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true).
			Width(7)

	MethodPATCH = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")). // Orange
			Bold(true).
			Width(7)

	MethodDELETE = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true).
			Width(7)

	// Input Styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Secondary).
				Padding(0, 1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginRight(1)

	// Status Styles
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Primary).
			Padding(0, 1)
)

func MethodStyle(method string) lipgloss.Style {
	switch method {
	case "GET":
		return MethodGET
	case "POST":
		return MethodPOST
	case "PUT":
		return MethodPUT
	case "PATCH":
		return MethodPATCH
	case "DELETE":
		return MethodDELETE
	default:
		return lipgloss.NewStyle().Width(7)
	}
}
