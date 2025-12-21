package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/internal/styles"
	"github.com/ksysoev/tapi/pkg/request"
)

func (m Model) handleRequestBuilderKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "j", "down":
		if len(m.inputs) > 0 {
			m.focusedInput = (m.focusedInput + 1) % len(m.inputs)
			return m, m.inputs[m.focusedInput].Focus()
		}
	case "shift+tab", "k", "up":
		if len(m.inputs) > 0 {
			m.focusedInput--
			if m.focusedInput < 0 {
				m.focusedInput = len(m.inputs) - 1
			}
			return m, m.inputs[m.focusedInput].Focus()
		}
	case "enter":
		if msg.Alt {
			return m, m.sendRequest()
		}
	case "ctrl+s":
		return m, m.sendRequest()
	}

	if len(m.inputs) > 0 {
		var cmd tea.Cmd
		m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) setupRequestBuilder() {
	op := m.getCurrentOperation()
	if op == nil {
		return
	}

	m.inputs = make([]textinput.Model, 0)

	for _, param := range op.Parameters {
		ti := textinput.New()
		ti.Placeholder = param.Name
		ti.CharLimit = 256
		ti.Width = 50
		ti.Prompt = fmt.Sprintf("%s (%s): ", param.Name, param.In)
		m.inputs = append(m.inputs, ti)
	}

	if op.RequestBody != nil && op.RequestBody.Required {
		ti := textinput.New()
		ti.Placeholder = "Request body (JSON)"
		ti.CharLimit = 2048
		ti.Width = 50
		ti.Prompt = "Body: "
		m.inputs = append(m.inputs, ti)
	}

	if len(m.inputs) > 0 {
		m.inputs[0].Focus()
		m.focusedInput = 0
	}
}

func (m Model) renderRequestBuilder() string {
	var b strings.Builder

	op := m.getCurrentOperation()
	if op == nil {
		return "No operation selected"
	}

	b.WriteString(styles.TitleStyle.Render("Request Builder"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtitleStyle.Render(fmt.Sprintf("%s %s", op.Method, m.getCurrentPath().Path)))
	b.WriteString("\n\n")

	if len(m.inputs) == 0 {
		b.WriteString(styles.SuccessStyle.Render("No parameters required"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("Press Ctrl+S to send request"))
	} else {
		for i, input := range m.inputs {
			if i == m.focusedInput {
				b.WriteString(styles.FocusedInputStyle.Render(input.View()))
			} else {
				b.WriteString(styles.InputStyle.Render(input.View()))
			}
			b.WriteString("\n\n")
		}
	}

	return b.String()
}

func (m Model) sendRequest() tea.Cmd {
	op := m.getCurrentOperation()
	path := m.getCurrentPath()
	
	if op == nil || path == nil {
		return nil
	}

	params := make(map[string]string)
	for i, input := range m.inputs {
		if i < len(op.Parameters) {
			params[op.Parameters[i].Name] = input.Value()
		}
	}

	var body string
	if op.RequestBody != nil && len(m.inputs) > len(op.Parameters) {
		body = m.inputs[len(m.inputs)-1].Value()
	}

	server := ""
	if len(m.spec.Servers) > 0 {
		server = m.spec.Servers[0].URL
	}

	return request.Send(server, path.Path, op.Method, params, body)
}
