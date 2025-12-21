package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kirill/tapi/internal/styles"
	"github.com/kirill/tapi/pkg/openapi"
	"github.com/kirill/tapi/pkg/request"
)

type view int

const (
	viewEndpoints view = iota
	viewOperationDetails
	viewRequestBuilder
	viewResponse
	viewHelp
)

type Model struct {
	spec              *openapi.Spec
	currentView       view
	endpointsList     []string
	selectedEndpoint  int
	selectedOperation int
	cursor            int
	width             int
	height            int
	viewport          viewport.Model
	inputs            []textinput.Model
	focusedInput      int
	lastResponse      string
	err               error
	showHelp          bool
}

func NewModel(spec *openapi.Spec) Model {
	endpoints := make([]string, 0)
	for _, path := range spec.Paths {
		for _, op := range path.Operations {
			endpoints = append(endpoints, fmt.Sprintf("%s %s", op.Method, path.Path))
		}
	}

	vp := viewport.New(80, 20)
	vp.Style = styles.PanelStyle

	return Model{
		spec:             spec,
		currentView:      viewEndpoints,
		endpointsList:    endpoints,
		selectedEndpoint: 0,
		viewport:         vp,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
	case request.ResponseMsg:
		m.lastResponse = msg.Body
		m.currentView = viewResponse
		m.viewport.SetContent(m.formatResponse(msg))
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		if m.currentView == viewEndpoints {
			return m, tea.Quit
		}
	case "?":
		m.showHelp = !m.showHelp
		return m, nil
	case "esc":
		if m.currentView != viewEndpoints {
			m.currentView = viewEndpoints
			m.showHelp = false
			return m, nil
		}
		return m, tea.Quit
	}

	// View-specific keys
	switch m.currentView {
	case viewEndpoints:
		return m.handleEndpointsKeys(msg)
	case viewOperationDetails:
		return m.handleOperationDetailsKeys(msg)
	case viewRequestBuilder:
		return m.handleRequestBuilderKeys(msg)
	case viewResponse:
		return m.handleResponseKeys(msg)
	}

	return m, nil
}

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
	}
	return m, nil
}

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

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	header := m.renderHeader()
	footer := m.renderFooter()

	switch m.currentView {
	case viewEndpoints:
		content = m.renderEndpoints()
	case viewOperationDetails:
		content = m.viewport.View()
	case viewRequestBuilder:
		content = m.renderRequestBuilder()
	case viewResponse:
		content = m.viewport.View()
	}

	if m.showHelp {
		content = m.renderHelp()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m Model) renderHeader() string {
	title := styles.TitleStyle.Render("🚀 TAPI - Terminal API Explorer")
	subtitle := styles.SubtitleStyle.Render(fmt.Sprintf("%s v%s", m.spec.Title, m.spec.Version))
	
	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle, "")
}

func (m Model) renderFooter() string {
	var keys string
	switch m.currentView {
	case viewEndpoints:
		keys = "j/k: navigate • enter: select • ?: help • q: quit"
	case viewOperationDetails:
		keys = "j/k: scroll • e: execute • h: back • ?: help • esc: exit"
	case viewRequestBuilder:
		keys = "tab: next field • ctrl+s: send request • h: back • esc: cancel"
	case viewResponse:
		keys = "j/k: scroll • h: back • esc: exit"
	}

	return styles.HelpStyle.Render(keys)
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

	return b.String()
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
	}

	return b.String()
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
	b.WriteString(resp.Body)

	return b.String()
}

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

func (m Model) getCurrentPath() *openapi.Path {
	if m.selectedEndpoint < 0 || m.selectedEndpoint >= len(m.endpointsList) {
		return nil
	}

	idx := 0
	for i, path := range m.spec.Paths {
		for range path.Operations {
			if idx == m.selectedEndpoint {
				return &m.spec.Paths[i]
			}
			idx++
		}
	}
	return nil
}

func (m Model) getCurrentOperation() *openapi.Operation {
	if m.selectedEndpoint < 0 || m.selectedEndpoint >= len(m.endpointsList) {
		return nil
	}

	idx := 0
	for _, path := range m.spec.Paths {
		for j := range path.Operations {
			if idx == m.selectedEndpoint {
				return &path.Operations[j]
			}
			idx++
		}
	}
	return nil
}
