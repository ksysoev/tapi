package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/pkg/openapi"
)

func TestHandleEndpointsKeysNavigation(t *testing.T) {
	tests := []struct {
		name             string
		key              string
		initialSelected  int
		expectedSelected int
	}{
		{"move down from first", "j", 0, 1},
		{"move down from middle", "down", 1, 2},
		{"move up from last", "k", 2, 1},
		{"move up from middle", "up", 1, 0},
		{"move up from first stays", "k", 0, 0},
		{"move down from last stays", "j", 2, 2},
		{"go to top", "g", 2, 0},
		{"go to bottom", "G", 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.initialSelected

			var msg tea.KeyMsg
			switch tt.key {
			case "j":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			case "k":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			case "g":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			case "G":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			case "down":
				msg = tea.KeyMsg{Type: tea.KeyDown}
			case "up":
				msg = tea.KeyMsg{Type: tea.KeyUp}
			}

			updatedModel, _ := model.handleEndpointsKeys(msg)
			m := updatedModel.(Model)

			if m.selectedEndpoint != tt.expectedSelected {
				t.Errorf("handleEndpointsKeys() selectedEndpoint = %d, want %d", 
					m.selectedEndpoint, tt.expectedSelected)
			}
		})
	}
}

func TestHandleEndpointsKeysEnter(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"enter key", "enter"},
		{"right arrow", "right"},
		{"l key", "l"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = 0

			var msg tea.KeyMsg
			switch tt.key {
			case "enter":
				msg = tea.KeyMsg{Type: tea.KeyEnter}
			case "right":
				msg = tea.KeyMsg{Type: tea.KeyRight}
			case "l":
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			}

			updatedModel, _ := model.handleEndpointsKeys(msg)
			m := updatedModel.(Model)

			if m.currentView != viewOperationDetails {
				t.Errorf("handleEndpointsKeys() currentView = %v, want %v", 
					m.currentView, viewOperationDetails)
			}
		})
	}
}

func TestRenderEndpoints(t *testing.T) {
	model := NewModel(createTestSpec())
	model.width = 100
	model.height = 50
	model.selectedEndpoint = 0

	rendered := model.renderEndpoints()

	if rendered == "" {
		t.Error("renderEndpoints() returned empty string")
	}

	if len(rendered) < 10 {
		t.Error("renderEndpoints() output seems too short")
	}
}

func TestRenderEndpointsWithScrolling(t *testing.T) {
	spec := createTestSpec()
	for i := 0; i < 20; i++ {
		spec.Paths = append(spec.Paths, openapi.Path{
			Path: "/extra/" + string(rune(i)),
			Operations: []openapi.Operation{
				{Method: "GET", Summary: "Extra endpoint"},
			},
		})
	}

	model := NewModel(spec)
	model.width = 100
	model.height = 30
	model.selectedEndpoint = 15

	rendered := model.renderEndpoints()

	if rendered == "" {
		t.Error("renderEndpoints() with scrolling returned empty string")
	}
}
