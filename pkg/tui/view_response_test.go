package tui

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/pkg/request"
)

func TestHandleResponseKeysNavigation(t *testing.T) {
	model := NewModel(createTestSpec())
	model.currentView = viewResponse
	model.viewport.SetContent("Response line 1\nResponse line 2\nResponse line 3")

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"move down with j", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}},
		{"move down with arrow", tea.KeyMsg{Type: tea.KeyDown}},
		{"move up with k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}},
		{"move up with arrow", tea.KeyMsg{Type: tea.KeyUp}},
		{"half page down", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}},
		{"half page up", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cmd := model.handleResponseKeys(tt.key)
			if cmd != nil {
				t.Errorf("handleResponseKeys() unexpected cmd")
			}
		})
	}
}

func TestHandleResponseKeysBack(t *testing.T) {
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
			model.currentView = viewResponse

			updatedModel, _ := model.handleResponseKeys(tt.key)
			m := updatedModel.(Model)

			if m.currentView != viewRequestBuilder {
				t.Errorf("handleResponseKeys() currentView = %v, want %v",
					m.currentView, viewRequestBuilder)
			}
		})
	}
}

func TestFormatResponseSuccess(t *testing.T) {
	model := NewModel(createTestSpec())

	resp := request.ResponseMsg{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":   []string{"application/json"},
			"Content-Length": []string{"123"},
		},
		Body: `{"result": "success", "data": [1, 2, 3]}`,
	}

	formatted := model.formatResponse(resp)

	if formatted == "" {
		t.Error("formatResponse() returned empty string")
	}

	expectedContains := []string{
		"Response: 200 OK",
		"Headers:",
		"Content-Type",
		"application/json",
		"Body:",
		`{"result": "success"`,
	}

	for _, want := range expectedContains {
		if !strings.Contains(formatted, want) {
			t.Errorf("formatResponse() missing %q in output", want)
		}
	}
}

func TestFormatResponseError(t *testing.T) {
	model := NewModel(createTestSpec())

	resp := request.ResponseMsg{
		Error: errors.New("connection timeout"),
	}

	formatted := model.formatResponse(resp)

	if !strings.Contains(formatted, "Error:") {
		t.Error("formatResponse() missing 'Error:' label")
	}

	if !strings.Contains(formatted, "connection timeout") {
		t.Error("formatResponse() missing error message")
	}
}

func TestFormatResponseMultipleHeaders(t *testing.T) {
	model := NewModel(createTestSpec())

	resp := request.ResponseMsg{
		StatusCode: 201,
		Status:     "Created",
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
			"X-Custom":     []string{"value1", "value2", "value3"},
		},
		Body: `{"id": 123}`,
	}

	formatted := model.formatResponse(resp)

	if !strings.Contains(formatted, "X-Custom") {
		t.Error("formatResponse() missing custom header")
	}

	if !strings.Contains(formatted, "value1") {
		t.Error("formatResponse() missing header value")
	}
}

func TestFormatResponseDifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		status     string
	}{
		{"200 OK", 200, "OK"},
		{"201 Created", 201, "Created"},
		{"400 Bad Request", 400, "Bad Request"},
		{"404 Not Found", 404, "Not Found"},
		{"500 Internal Server Error", 500, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())

			resp := request.ResponseMsg{
				StatusCode: tt.statusCode,
				Status:     tt.status,
				Headers:    http.Header{},
				Body:       `{}`,
			}

			formatted := model.formatResponse(resp)

			expectedStatus := "Response: " + tt.name
			if !strings.Contains(formatted, expectedStatus) {
				t.Errorf("formatResponse() missing %q in output", expectedStatus)
			}
		})
	}
}

func TestFormatResponseEmptyBody(t *testing.T) {
	model := NewModel(createTestSpec())

	resp := request.ResponseMsg{
		StatusCode: 204,
		Status:     "No Content",
		Headers:    http.Header{},
		Body:       "",
	}

	formatted := model.formatResponse(resp)

	if !strings.Contains(formatted, "Body:") {
		t.Error("formatResponse() should still show Body label")
	}

	if !strings.Contains(formatted, "Response: 204 No Content") {
		t.Error("formatResponse() missing status line")
	}
}
