package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksysoev/tapi/internal/styles"
)

func (m Model) handleOperationDetailsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.viewport.LineDown(1)
	case "k", "up":
		m.viewport.LineUp(1)
	case "d":
		m.viewport.HalfViewDown()
	case "u":
		m.viewport.HalfViewUp()
	case "e", "enter":
		m.currentView = viewRequestBuilder
		m.setupRequestBuilder()
	case "h", "left":
		m.currentView = viewEndpoints
	}
	return m, nil
}

func (m Model) getOperationDetails() string {
	op := m.getCurrentOperation()
	if op == nil {
		return "No operation selected"
	}

	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render(fmt.Sprintf("%s %s", op.Method, m.getCurrentPath().Path)))
	b.WriteString("\n\n")

	if op.Summary != "" {
		b.WriteString(styles.LabelStyle.Render("Summary: "))
		b.WriteString(op.Summary)
		b.WriteString("\n\n")
	}

	if op.Description != "" {
		b.WriteString(styles.LabelStyle.Render("Description:"))
		b.WriteString("\n")
		b.WriteString(op.Description)
		b.WriteString("\n\n")
	}

	if len(op.Parameters) > 0 {
		b.WriteString(styles.LabelStyle.Render("Parameters:"))
		b.WriteString("\n")
		for _, param := range op.Parameters {
			required := ""
			if param.Required {
				required = lipgloss.NewStyle().Foreground(styles.Danger).Render(" *")
			}
			b.WriteString(fmt.Sprintf("  • %s (%s)%s - %s\n", param.Name, param.In, required, param.Description))
		}
		b.WriteString("\n")
	}

	if op.RequestBody != nil {
		b.WriteString(styles.LabelStyle.Render("Request Body:"))
		b.WriteString("\n")
		for contentType := range op.RequestBody.Content {
			b.WriteString(fmt.Sprintf("  • %s\n", contentType))
		}
		b.WriteString("\n")
	}

	if len(op.Responses) > 0 {
		b.WriteString(styles.LabelStyle.Render("Responses:"))
		b.WriteString("\n")
		for status, resp := range op.Responses {
			b.WriteString(fmt.Sprintf("  • %s - %s\n", status, resp.Description))
		}
		b.WriteString("\n")
	}

	// Help to fix issue that content is not possible to scroll down fully
	b.WriteString("\n\n\n\n")

	return b.String()
}
