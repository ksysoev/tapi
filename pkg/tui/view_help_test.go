package tui

import (
	"strings"
	"testing"
)

func TestRenderHelp(t *testing.T) {
	model := NewModel(createTestSpec())

	help := model.renderHelp()

	if help == "" {
		t.Error("renderHelp() returned empty string")
	}

	expectedSections := []string{
		"TAPI - Terminal API Explorer",
		"Navigation:",
		"Actions:",
		"General:",
		"j, ↓",
		"k, ↑",
		"h, ←",
		"l, →",
		"Enter",
		"Ctrl+S",
		"?",
		"Esc",
		"q",
		"Ctrl+C",
	}

	for _, section := range expectedSections {
		if !strings.Contains(help, section) {
			t.Errorf("renderHelp() missing section or key: %q", section)
		}
	}
}

func TestRenderHelpContainsAllKeyBindings(t *testing.T) {
	model := NewModel(createTestSpec())
	help := model.renderHelp()

	keyBindings := []string{
		"Move down",
		"Move up",
		"Go back",
		"Move forward",
		"Go to top",
		"Go to bottom",
		"Scroll half page down",
		"Scroll half page up",
		"Select",
		"Confirm",
		"Execute API request",
		"Send request",
		"Next input field",
		"Previous input field",
		"Toggle help",
		"Quit",
		"Force quit",
	}

	for _, binding := range keyBindings {
		if !strings.Contains(help, binding) {
			t.Errorf("renderHelp() missing key binding description: %q", binding)
		}
	}
}
