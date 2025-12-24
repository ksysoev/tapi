package main

import (
	"testing"
)

func TestRunApp(t *testing.T) {
	// Test that runApp can be called without panicking
	// We can't test full execution as it requires command line args
	// This at least verifies the function structure is valid
	
	// Just verify the command structure is valid
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("runApp panicked: %v", r)
		}
	}()
	
	// The function will likely return 1 (error) without proper args,
	// but we're just checking it doesn't panic
	_ = runApp()
}

func TestVersionAndName(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
	
	if name == "" {
		t.Error("name should not be empty")
	}
	
	if name != "tapi" {
		t.Errorf("Expected name to be 'tapi', got '%s'", name)
	}
}
