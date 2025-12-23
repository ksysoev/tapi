package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/internal/styles"
)

func (m Model) handleEndpointsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selectedEndpoint < len(m.endpointsList)-1 {
			m.selectedEndpoint++
		}
	case "k", "up":
		if m.selectedEndpoint > 0 {
			m.selectedEndpoint--
		}
	case "g":
		m.selectedEndpoint = 0
	case "G":
		m.selectedEndpoint = len(m.endpointsList) - 1
	case "enter", "l", "right":
		m.currentView = viewOperationDetails
		m.viewport.SetContent(m.getOperationDetails())
		m.viewport.GotoTop()
	}
	return m, nil
}

func (m Model) renderEndpoints() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Endpoints"))
	b.WriteString("\n\n")

	start := 0
	end := len(m.endpointsList)

	maxVisible := m.height - 15
	if end-start > maxVisible {
		if m.selectedEndpoint > maxVisible/2 {
			start = m.selectedEndpoint - maxVisible/2
		}
		if end-start > maxVisible {
			end = start + maxVisible
		}
	}

	for i := start; i < end && i < len(m.endpointsList); i++ {
		parts := strings.Fields(m.endpointsList[i])
		if len(parts) >= 2 {
			method := parts[0]
			path := strings.Join(parts[1:], " ")

			line := fmt.Sprintf("%s %s",
				styles.MethodStyle(method).Render(method),
				path,
			)

			if i == m.selectedEndpoint {
				b.WriteString(styles.SelectedItemStyle.Render("▶ " + line))
			} else {
				b.WriteString(styles.ItemStyle.Render(line))
			}
			b.WriteString("\n")
		}
	}

	// Help to fix issue that content is not possible to scroll down fully
	b.WriteString("\n\n\n\n")

	return b.String()
}
