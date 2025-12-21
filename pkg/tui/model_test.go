package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/pkg/openapi"
	"github.com/ksysoev/tapi/pkg/request"
)

func createTestSpec() *openapi.Spec {
	return &openapi.Spec{
		Title:   "Test API",
		Version: "1.0.0",
		Servers: []openapi.Server{
			{URL: "https://api.example.com", Description: "Test server"},
		},
		Paths: []openapi.Path{
			{
				Path: "/users",
				Operations: []openapi.Operation{
					{
						Method:      "GET",
						Summary:     "Get users",
						Description: "Retrieve all users",
						Parameters: []openapi.Parameter{
							{Name: "limit", In: "query", Description: "Limit results", Required: false},
						},
						Responses: map[string]openapi.Response{
							"200": {Description: "Success"},
						},
					},
					{
						Method:      "POST",
						Summary:     "Create user",
						Description: "Create a new user",
						RequestBody: &openapi.RequestBody{
							Required: true,
							Content: map[string]openapi.MediaType{
								"application/json": {},
							},
						},
						Responses: map[string]openapi.Response{
							"201": {Description: "Created"},
						},
					},
				},
			},
			{
				Path: "/users/{id}",
				Operations: []openapi.Operation{
					{
						Method:      "GET",
						Summary:     "Get user by ID",
						Description: "Retrieve a single user",
						Parameters: []openapi.Parameter{
							{Name: "id", In: "path", Description: "User ID", Required: true},
						},
						Responses: map[string]openapi.Response{
							"200": {Description: "Success"},
							"404": {Description: "Not found"},
						},
					},
				},
			},
		},
	}
}

func TestNewModel(t *testing.T) {
	spec := createTestSpec()
	model := NewModel(spec)

	if model.spec != spec {
		t.Error("NewModel() spec not set correctly")
	}

	if model.currentView != viewEndpoints {
		t.Errorf("NewModel() currentView = %v, want %v", model.currentView, viewEndpoints)
	}

	expectedEndpoints := 3
	if len(model.endpointsList) != expectedEndpoints {
		t.Errorf("NewModel() endpointsList length = %d, want %d", len(model.endpointsList), expectedEndpoints)
	}

	if model.selectedEndpoint != 0 {
		t.Errorf("NewModel() selectedEndpoint = %d, want 0", model.selectedEndpoint)
	}

	expected := "GET /users"
	if model.endpointsList[0] != expected {
		t.Errorf("NewModel() first endpoint = %q, want %q", model.endpointsList[0], expected)
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel(createTestSpec())
	cmd := model.Init()

	if cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	model := NewModel(createTestSpec())
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Update() width = %d, want 100", m.width)
	}

	if m.height != 50 {
		t.Errorf("Update() height = %d, want 50", m.height)
	}

	expectedViewportWidth := 96
	if m.viewport.Width != expectedViewportWidth {
		t.Errorf("Update() viewport width = %d, want %d", m.viewport.Width, expectedViewportWidth)
	}

	expectedViewportHeight := 40
	if m.viewport.Height != expectedViewportHeight {
		t.Errorf("Update() viewport height = %d, want %d", m.viewport.Height, expectedViewportHeight)
	}
}

func TestModelUpdateResponseMsg(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewRequestBuilder

	responseMsg := request.ResponseMsg{
		StatusCode: 200,
		Status:     "OK",
		Body:       `{"result": "success"}`,
		Headers:    map[string][]string{"Content-Type": {"application/json"}},
	}

	updatedModel, _ := model.Update(responseMsg)
	m := updatedModel.(Model)

	if m.currentView != viewResponse {
		t.Errorf("Update() currentView = %v, want %v", m.currentView, viewResponse)
	}

	if m.lastResponse != responseMsg.Body {
		t.Errorf("Update() lastResponse = %q, want %q", m.lastResponse, responseMsg.Body)
	}
}

func TestHandleKeyPressQuit(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		view        view
		shouldQuit  bool
		description string
	}{
		{"quit from endpoints", "q", viewEndpoints, true, "q should quit from endpoints view"},
		{"quit from details", "q", viewOperationDetails, false, "q should not quit from details view"},
		{"ctrl+c from endpoints", "ctrl+c", viewEndpoints, true, "ctrl+c should quit from endpoints"},
		{"ctrl+c from details", "ctrl+c", viewOperationDetails, false, "ctrl+c should not quit from details view"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.currentView = tt.view

			msg := tea.KeyMsg{Type: tea.KeyRunes}
			switch tt.key {
			case "q":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			case "ctrl+c":
				msg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := model.handleKeyPress(msg)

			gotQuit := cmd != nil && cmd() == tea.Quit()
			if gotQuit != tt.shouldQuit {
				t.Errorf("%s: shouldQuit = %v, want %v", tt.description, gotQuit, tt.shouldQuit)
			}
		})
	}
}

func TestHandleKeyPressHelp(t *testing.T) {
	model := NewModel(createTestSpec())
	model.showHelp = false

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, _ := model.handleKeyPress(msg)
	m := updatedModel.(Model)

	if !m.showHelp {
		t.Error("handleKeyPress() should toggle help on")
	}

	updatedModel, _ = m.handleKeyPress(msg)
	m = updatedModel.(Model)

	if m.showHelp {
		t.Error("handleKeyPress() should toggle help off")
	}
}

func TestHandleKeyPressEscape(t *testing.T) {
	tests := []struct {
		name         string
		initialView  view
		expectedView view
		shouldQuit   bool
	}{
		{"escape from details", viewOperationDetails, viewEndpoints, false},
		{"escape from request builder", viewRequestBuilder, viewEndpoints, false},
		{"escape from response", viewResponse, viewEndpoints, false},
		{"escape from endpoints", viewEndpoints, viewEndpoints, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.currentView = tt.initialView
			model.showHelp = true

			msg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, cmd := model.handleKeyPress(msg)
			m := updatedModel.(Model)

			if m.currentView != tt.expectedView {
				t.Errorf("handleKeyPress() currentView = %v, want %v", m.currentView, tt.expectedView)
			}

			if tt.initialView != viewEndpoints && m.showHelp {
				t.Error("handleKeyPress() should hide help on escape from non-endpoints view")
			}

			gotQuit := cmd != nil && cmd() == tea.Quit()
			if gotQuit != tt.shouldQuit {
				t.Errorf("handleKeyPress() shouldQuit = %v, want %v", gotQuit, tt.shouldQuit)
			}
		})
	}
}

func TestView(t *testing.T) {
	model := NewModel(createTestSpec())

	view := model.View()
	if view != "Loading..." {
		t.Errorf("View() before window size = %q, want 'Loading...'", view)
	}

	model.width = 100
	model.height = 50

	view = model.View()
	if view == "" {
		t.Error("View() returned empty string")
	}

	if view == "Loading..." {
		t.Error("View() still shows loading after size set")
	}
}

func TestRenderHeader(t *testing.T) {
	model := NewModel(createTestSpec())
	header := model.renderHeader()

	if header == "" {
		t.Error("renderHeader() returned empty string")
	}
}

func TestRenderFooter(t *testing.T) {
	tests := []struct {
		name string
		view view
	}{
		{"endpoints footer", viewEndpoints},
		{"details footer", viewOperationDetails},
		{"request builder footer", viewRequestBuilder},
		{"response footer", viewResponse},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.currentView = tt.view

			footer := model.renderFooter()
			if footer == "" {
				t.Error("renderFooter() returned empty string")
			}
		})
	}
}
