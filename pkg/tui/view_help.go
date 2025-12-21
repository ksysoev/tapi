package tui

import "github.com/ksysoev/tapi/internal/styles"

func (m Model) renderHelp() string {
	help := `
TAPI - Terminal API Explorer - Keyboard Shortcuts

Navigation:
  j, ↓          Move down
  k, ↑          Move up
  h, ←          Go back / Move left
  l, →          Move forward / Move right
  g             Go to top
  G             Go to bottom
  d             Scroll half page down
  u             Scroll half page up

Actions:
  Enter         Select / Confirm
  e             Execute API request
  Ctrl+S        Send request
  Tab           Next input field
  Shift+Tab     Previous input field

General:
  ?             Toggle help
  Esc           Go back / Cancel
  q             Quit (from main view)
  Ctrl+C        Force quit
`

	return styles.PanelStyle.Render(help)
}
