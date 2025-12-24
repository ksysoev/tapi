package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ksysoev/tapi/internal/styles"
	"github.com/ksysoev/tapi/pkg/openapi"
	"github.com/ksysoev/tapi/pkg/request"
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
	width             int
	height            int
	viewport          viewport.Model
	inputs            []textinput.Model
	focusedInput      int
	lastResponse      string
	showHelp          bool
}

func NewModel(spec *openapi.Spec) Model {
	endpoints := make([]string, 0)
	for _, path := range spec.Paths {
		for _, op := range path.Operations {
			endpoints = append(endpoints, fmt.Sprintf("%s %s", op.Method, path.Path))
		}
	}

	// Sort endpoints by path first, then by method
	sort.Slice(endpoints, func(i, j int) bool {
		partsI := strings.Fields(endpoints[i])
		partsJ := strings.Fields(endpoints[j])
		
		if len(partsI) < 2 || len(partsJ) < 2 {
			return endpoints[i] < endpoints[j]
		}
		
		pathI := strings.Join(partsI[1:], " ")
		pathJ := strings.Join(partsJ[1:], " ")
		
		if pathI == pathJ {
			return partsI[0] < partsJ[0]
		}
		return pathI < pathJ
	})

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
		m.viewport.YOffset = 0
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
