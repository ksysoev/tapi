package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/pkg/openapi"
)

func TestHandleRequestBuilderKeysNavigation(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewRequestBuilder
	model.selectedEndpoint = 1
	model.setupRequestBuilder()

	if len(model.inputs) < 2 {
		t.Skip("Need at least 2 inputs for navigation test")
	}

	tests := []struct {
		name             string
		key              tea.KeyMsg
		initialFocus     int
		expectedFocus    int
	}{
		{"tab next", tea.KeyMsg{Type: tea.KeyTab}, 0, 1},
		{"shift+tab prev", tea.KeyMsg{Type: tea.KeyShiftTab}, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testModel := model
			testModel.focusedInput = tt.initialFocus

			updatedModel, cmd := testModel.handleRequestBuilderKeys(tt.key)
			m := updatedModel.(Model)

			if m.focusedInput != tt.expectedFocus {
				t.Errorf("handleRequestBuilderKeys() focusedInput = %d, want %d",
					m.focusedInput, tt.expectedFocus)
			}

			if cmd == nil {
				t.Error("handleRequestBuilderKeys() should return focus cmd")
			}
		})
	}
}

func TestHandleRequestBuilderKeysNavigationWraparound(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewRequestBuilder
	model.selectedEndpoint = 0
	model.setupRequestBuilder()

	numInputs := len(model.inputs)
	model.focusedInput = numInputs - 1

	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.handleRequestBuilderKeys(msg)
	m := updatedModel.(Model)

	if m.focusedInput != 0 {
		t.Errorf("handleRequestBuilderKeys() tab wraparound focusedInput = %d, want 0", m.focusedInput)
	}

	m.focusedInput = 0
	msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	updatedModel, _ = m.handleRequestBuilderKeys(msg)
	m = updatedModel.(Model)

	if m.focusedInput != numInputs-1 {
		t.Errorf("handleRequestBuilderKeys() shift+tab wraparound focusedInput = %d, want %d",
			m.focusedInput, numInputs-1)
	}
}

func TestHandleRequestBuilderKeysSend(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewRequestBuilder
	model.selectedEndpoint = 0
	model.setupRequestBuilder()

	msg := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, cmd := model.handleRequestBuilderKeys(msg)

	if cmd == nil {
		t.Error("handleRequestBuilderKeys() with ctrl+s should return send cmd")
	}
}

func TestSetupRequestBuilder(t *testing.T) {
	tests := []struct {
		name             string
		selectedEndpoint int
		expectedInputs   int
	}{
		{"GET with one param", 0, 1},
		{"POST with request body", 1, 1},
		{"GET with path param", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.selectedEndpoint

			model.setupRequestBuilder()

			if len(model.inputs) != tt.expectedInputs {
				t.Errorf("setupRequestBuilder() inputs count = %d, want %d",
					len(model.inputs), tt.expectedInputs)
			}

			if len(model.inputs) > 0 && model.focusedInput != 0 {
				t.Errorf("setupRequestBuilder() focusedInput = %d, want 0", model.focusedInput)
			}
		})
	}
}

func TestSetupRequestBuilderInvalidEndpoint(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 100

	model.setupRequestBuilder()

	if len(model.inputs) != 0 {
		t.Errorf("setupRequestBuilder() with invalid endpoint should have 0 inputs, got %d", len(model.inputs))
	}
}

func TestRenderRequestBuilder(t *testing.T) {
	tests := []struct {
		name             string
		selectedEndpoint int
		wantContains     []string
	}{
		{
			name:             "GET with param",
			selectedEndpoint: 0,
			wantContains:     []string{"Request Builder", "GET /users"},
		},
		{
			name:             "POST with body",
			selectedEndpoint: 1,
			wantContains:     []string{"Request Builder", "POST /users"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.selectedEndpoint
			model.setupRequestBuilder()

			rendered := model.renderRequestBuilder()

			if rendered == "" {
				t.Error("renderRequestBuilder() returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(rendered, want) {
					t.Errorf("renderRequestBuilder() missing %q in output", want)
				}
			}
		})
	}
}

func TestRenderRequestBuilderNoParams(t *testing.T) {
	spec := createTestSpec()
	spec.Paths[0].Operations[0].Parameters = []openapi.Parameter{}

	model := NewModel(spec)
	model.selectedEndpoint = 0
	model.setupRequestBuilder()

	rendered := model.renderRequestBuilder()

	if !strings.Contains(rendered, "No parameters required") {
		t.Error("renderRequestBuilder() should show 'No parameters required' message")
	}

	if !strings.Contains(rendered, "Ctrl+S") {
		t.Error("renderRequestBuilder() should show send instruction")
	}
}

func TestRenderRequestBuilderInvalidEndpoint(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 100

	rendered := model.renderRequestBuilder()

	if rendered != "No operation selected" {
		t.Errorf("renderRequestBuilder() = %q, want 'No operation selected'", rendered)
	}
}

func TestSendRequest(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 0
	model.setupRequestBuilder()

	model.inputs[0].SetValue("10")

	cmd := model.sendRequest()

	if cmd == nil {
		t.Error("sendRequest() should return command")
	}
}

func TestSendRequestWithBody(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 1
	model.setupRequestBuilder()

	model.inputs[0].SetValue(`{"name":"John"}`)

	cmd := model.sendRequest()

	if cmd == nil {
		t.Error("sendRequest() should return command")
	}
}

func TestSendRequestInvalidEndpoint(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 100

	cmd := model.sendRequest()

	if cmd != nil {
		t.Error("sendRequest() should return nil for invalid endpoint")
	}
}

func TestSendRequestNoServer(t *testing.T) {
	spec := createTestSpec()
	spec.Servers = []openapi.Server{}

	model := NewModel(spec)
	model.selectedEndpoint = 0
	model.setupRequestBuilder()

	cmd := model.sendRequest()

	if cmd == nil {
		t.Error("sendRequest() should return command even without server")
	}
}
