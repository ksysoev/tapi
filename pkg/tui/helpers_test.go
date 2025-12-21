package tui

import (
	"testing"
)

func TestGetCurrentPath(t *testing.T) {
	tests := []struct {
		name             string
		selectedEndpoint int
		wantNil          bool
		expectedPath     string
	}{
		{"first endpoint", 0, false, "/users"},
		{"second endpoint", 1, false, "/users"},
		{"third endpoint", 2, false, "/users/{id}"},
		{"invalid negative", -1, true, ""},
		{"invalid too large", 100, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.selectedEndpoint

			path := model.getCurrentPath()

			if tt.wantNil {
				if path != nil {
					t.Errorf("getCurrentPath() = %v, want nil", path)
				}
			} else {
				if path == nil {
					t.Fatal("getCurrentPath() returned nil, want non-nil")
				}
				if path.Path != tt.expectedPath {
					t.Errorf("getCurrentPath() path = %q, want %q", path.Path, tt.expectedPath)
				}
			}
		})
	}
}

func TestGetCurrentOperation(t *testing.T) {
	tests := []struct {
		name             string
		selectedEndpoint int
		wantNil          bool
		expectedMethod   string
		expectedSummary  string
	}{
		{"first operation GET /users", 0, false, "GET", "Get users"},
		{"second operation POST /users", 1, false, "POST", "Create user"},
		{"third operation GET /users/{id}", 2, false, "GET", "Get user by ID"},
		{"invalid negative", -1, true, "", ""},
		{"invalid too large", 100, true, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel(createTestSpec())
			model.selectedEndpoint = tt.selectedEndpoint

			op := model.getCurrentOperation()

			if tt.wantNil {
				if op != nil {
					t.Errorf("getCurrentOperation() = %v, want nil", op)
				}
			} else {
				if op == nil {
					t.Fatal("getCurrentOperation() returned nil, want non-nil")
				}
				if op.Method != tt.expectedMethod {
					t.Errorf("getCurrentOperation() method = %q, want %q", op.Method, tt.expectedMethod)
				}
				if op.Summary != tt.expectedSummary {
					t.Errorf("getCurrentOperation() summary = %q, want %q", op.Summary, tt.expectedSummary)
				}
			}
		})
	}
}
