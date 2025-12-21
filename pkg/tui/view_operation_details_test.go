package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHandleOperationDetailsKeysNavigation(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewOperationDetails
	model.viewport.SetContent("Line 1\nLine 2\nLine 3\nLine 4\nLine 5")

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"move down", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}},
		{"move up", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}},
		{"arrow down", tea.KeyMsg{Type: tea.KeyDown}},
		{"arrow up", tea.KeyMsg{Type: tea.KeyUp}},
		{"half page down", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}},
		{"half page up", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cmd := model.handleOperationDetailsKeys(tt.key)
			if cmd != nil {
				t.Errorf("handleOperationDetailsKeys() unexpected cmd")
			}
		})
	}
}

func TestHandleOperationDetailsKeysExecute(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"e key", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}},
		{"enter key", tea.KeyMsg{Type: tea.KeyEnter}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.currentView = viewOperationDetails
			model.selectedEndpoint = 0

			updatedModel, _ := model.handleOperationDetailsKeys(tt.key)
			m := updatedModel.(Model)

			if m.currentView != viewRequestBuilder {
				t.Errorf("handleOperationDetailsKeys() currentView = %v, want %v",
					m.currentView, viewRequestBuilder)
			}

			if len(m.inputs) == 0 {
				t.Error("handleOperationDetailsKeys() did not setup inputs")
			}
		})
	}
}

func TestHandleOperationDetailsKeysBack(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"h key", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}},
		{"left arrow", tea.KeyMsg{Type: tea.KeyLeft}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.currentView = viewOperationDetails

			updatedModel, _ := model.handleOperationDetailsKeys(tt.key)
			m := updatedModel.(Model)

			if m.currentView != viewEndpoints {
				t.Errorf("handleOperationDetailsKeys() currentView = %v, want %v",
					m.currentView, viewEndpoints)
			}
		})
	}
}

func TestGetOperationDetails(t *testing.T) {
	tests := []struct {
		name             string
		selectedEndpoint int
		wantContains     []string
	}{
		{
			name:             "GET /users details",
			selectedEndpoint: 0,
			wantContains:     []string{"GET /users", "Get users", "Retrieve all users", "Parameters:", "limit"},
		},
		{
			name:             "POST /users details",
			selectedEndpoint: 1,
			wantContains:     []string{"POST /users", "Create user", "Create a new user", "Request Body:", "application/json"},
		},
		{
			name:             "GET /users/{id} details",
			selectedEndpoint: 2,
			wantContains:     []string{"GET /users/{id}", "Get user by ID", "Retrieve a single user", "id", "path"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.selectedEndpoint

			details := model.getOperationDetails()

			if details == "" {
				t.Error("getOperationDetails() returned empty string")
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(details, want) {
					t.Errorf("getOperationDetails() missing %q in output", want)
				}
			}
		})
	}
}

func TestGetOperationDetailsNoSelection(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 100

	details := model.getOperationDetails()

	if details != "No operation selected" {
		t.Errorf("getOperationDetails() = %q, want 'No operation selected'", details)
	}
}

func TestGetOperationDetailsWithResponses(t *testing.T) {
	model := NewModel(createTestSpec())
	model.selectedEndpoint = 2

	details := model.getOperationDetails()

	if !strings.Contains(details, "Responses:") {
		t.Error("getOperationDetails() missing Responses section")
	}

	if !strings.Contains(details, "200") {
		t.Error("getOperationDetails() missing 200 response")
	}

	if !strings.Contains(details, "404") {
		t.Error("getOperationDetails() missing 404 response")
	}
}
