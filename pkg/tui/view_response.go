package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/internal/styles"
	"github.com/ksysoev/tapi/pkg/formatter"
	"github.com/ksysoev/tapi/pkg/request"
)

func (m Model) handleResponseKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.viewport.LineDown(1)
	case "k", "up":
		m.viewport.LineUp(1)
	case "d":
		m.viewport.HalfViewDown()
	case "u":
		m.viewport.HalfViewUp()
	case "h", "left":
		m.currentView = viewRequestBuilder
	}
	return m, nil
}

func (m Model) formatResponse(resp request.ResponseMsg) string {
	var b strings.Builder

	if resp.Error != nil {
		b.WriteString(styles.ErrorStyle.Render("Error: "))
		b.WriteString(resp.Error.Error())
		return b.String()
	}

	b.WriteString(styles.SuccessStyle.Render(fmt.Sprintf("Response: %d %s", resp.StatusCode, resp.Status)))
	b.WriteString("\n\n")

	b.WriteString(styles.LabelStyle.Render("Headers:"))
	b.WriteString("\n")
	for key, values := range resp.Headers {
		b.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
	}

	b.WriteString("\n")
	b.WriteString(styles.LabelStyle.Render("Body:"))
	b.WriteString("\n")

	contentType := resp.Headers.Get("Content-Type")
	formattedBody := formatter.DetectAndFormat(resp.Body, contentType)
	b.WriteString(formattedBody)

	// Help to fix issue that content is not possible to scroll down fully
	b.WriteString("\n\n\n\n")

	return b.String()
}
