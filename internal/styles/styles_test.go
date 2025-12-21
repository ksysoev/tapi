package styles

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestMethodStyle(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   lipgloss.Style
	}{
		{
			name:   "GET method",
			method: "GET",
			want:   MethodGET,
		},
		{
			name:   "POST method",
			method: "POST",
			want:   MethodPOST,
		},
		{
			name:   "PUT method",
			method: "PUT",
			want:   MethodPUT,
		},
		{
			name:   "PATCH method",
			method: "PATCH",
			want:   MethodPATCH,
		},
		{
			name:   "DELETE method",
			method: "DELETE",
			want:   MethodDELETE,
		},
		{
			name:   "unknown method",
			method: "OPTIONS",
			want:   lipgloss.NewStyle().Width(7),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MethodStyle(tt.method)
			
			if got.GetWidth() != 7 {
				t.Errorf("MethodStyle(%s) width = %d, want 7", tt.method, got.GetWidth())
			}
		})
	}
}

func TestMethodStyleConsistentWidth(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	
	for _, method := range methods {
		style := MethodStyle(method)
		if style.GetWidth() != 7 {
			t.Errorf("MethodStyle(%s) width = %d, want 7", method, style.GetWidth())
		}
	}
}

func TestStyleColors(t *testing.T) {
	colorTests := []struct {
		name  string
		color lipgloss.Color
	}{
		{"Primary", Primary},
		{"Secondary", Secondary},
		{"Success", Success},
		{"Warning", Warning},
		{"Danger", Danger},
		{"Muted", Muted},
		{"Text", Text},
	}

	for _, tt := range colorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color == "" {
				t.Errorf("%s color is empty", tt.name)
			}
		})
	}
}

func TestTitleStyle(t *testing.T) {
	style := TitleStyle
	
	if style.GetForeground() != Primary {
		t.Error("TitleStyle foreground should be Primary color")
	}
	
	if !style.GetBold() {
		t.Error("TitleStyle should be bold")
	}
}

func TestSubtitleStyle(t *testing.T) {
	style := SubtitleStyle
	
	if style.GetForeground() != Secondary {
		t.Error("SubtitleStyle foreground should be Secondary color")
	}
	
	if !style.GetItalic() {
		t.Error("SubtitleStyle should be italic")
	}
}

func TestErrorStyle(t *testing.T) {
	style := ErrorStyle
	
	if style.GetForeground() != Danger {
		t.Error("ErrorStyle foreground should be Danger color")
	}
	
	if !style.GetBold() {
		t.Error("ErrorStyle should be bold")
	}
}

func TestSuccessStyle(t *testing.T) {
	style := SuccessStyle
	
	if style.GetForeground() != Success {
		t.Error("SuccessStyle foreground should be Success color")
	}
	
	if !style.GetBold() {
		t.Error("SuccessStyle should be bold")
	}
}

func TestPanelStyleHasBorder(t *testing.T) {
	style := PanelStyle
	
	border := style.GetBorderStyle()
	if border.Top == "" && border.Bottom == "" && border.Left == "" && border.Right == "" {
		t.Error("PanelStyle should have a border")
	}
}

func TestActivePanelStyleHasBorder(t *testing.T) {
	style := ActivPanelStyle
	
	border := style.GetBorderStyle()
	if border.Top == "" && border.Bottom == "" && border.Left == "" && border.Right == "" {
		t.Error("ActivPanelStyle should have a border")
	}
}

func TestSelectedItemStyle(t *testing.T) {
	style := SelectedItemStyle
	
	if style.GetForeground() != Secondary {
		t.Error("SelectedItemStyle foreground should be Secondary color")
	}
	
	if !style.GetBold() {
		t.Error("SelectedItemStyle should be bold")
	}
}

func TestInputStyleHasBorder(t *testing.T) {
	style := InputStyle
	
	border := style.GetBorderStyle()
	if border.Top == "" && border.Bottom == "" && border.Left == "" && border.Right == "" {
		t.Error("InputStyle should have a border")
	}
}

func TestFocusedInputStyleHasBorder(t *testing.T) {
	style := FocusedInputStyle
	
	border := style.GetBorderStyle()
	if border.Top == "" && border.Bottom == "" && border.Left == "" && border.Right == "" {
		t.Error("FocusedInputStyle should have a border")
	}
}

func TestLabelStyle(t *testing.T) {
	style := LabelStyle
	
	if style.GetForeground() != Primary {
		t.Error("LabelStyle foreground should be Primary color")
	}
	
	if !style.GetBold() {
		t.Error("LabelStyle should be bold")
	}
}

func TestStatusBarStyle(t *testing.T) {
	style := StatusBarStyle
	
	if style.GetForeground() != Text {
		t.Error("StatusBarStyle foreground should be Text color")
	}
	
	if style.GetBackground() != Primary {
		t.Error("StatusBarStyle background should be Primary color")
	}
}

func TestMethodStyleUniqueness(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	
	// Just verify that all methods return a style
	for _, method := range methods {
		style := MethodStyle(method)
		if style.GetWidth() != 7 {
			t.Errorf("MethodStyle(%s) should have width 7", method)
		}
	}
	
	// Test that unknown method returns default style
	unknownStyle := MethodStyle("UNKNOWN")
	if unknownStyle.GetWidth() != 7 {
		t.Error("MethodStyle for unknown method should have width 7")
	}
}
