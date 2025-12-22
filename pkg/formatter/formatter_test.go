package formatter

import (
	"strings"
	"testing"
)

func TestJSONFormatterCanHandle(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name        string
		content     string
		contentType string
		want        bool
	}{
		{
			name:        "JSON content type",
			content:     `{"key": "value"}`,
			contentType: "application/json",
			want:        true,
		},
		{
			name:        "JSON object",
			content:     `{"key": "value"}`,
			contentType: "",
			want:        true,
		},
		{
			name:        "JSON array",
			content:     `[1, 2, 3]`,
			contentType: "",
			want:        true,
		},
		{
			name:        "Empty content",
			content:     "",
			contentType: "",
			want:        false,
		},
		{
			name:        "Plain text",
			content:     "Hello world",
			contentType: "text/plain",
			want:        false,
		},
		{
			name:        "JSON with whitespace",
			content:     "  \n  {\"key\": \"value\"}  ",
			contentType: "",
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatter.CanHandle(tt.content, tt.contentType); got != tt.want {
				t.Errorf("CanHandle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONFormatterFormat(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name     string
		content  string
		contains []string
	}{
		{
			name:    "Simple JSON object",
			content: `{"name":"John","age":30}`,
			contains: []string{
				"name",
				"John",
				"age",
				"30",
			},
		},
		{
			name:    "Nested JSON",
			content: `{"user":{"name":"Jane","active":true}}`,
			contains: []string{
				"user",
				"name",
				"Jane",
				"active",
				"true",
			},
		},
		{
			name:    "JSON with array",
			content: `{"items":[1,2,3]}`,
			contains: []string{
				"items",
				"1",
				"2",
				"3",
			},
		},
		{
			name:    "JSON with null",
			content: `{"value":null}`,
			contains: []string{
				"value",
				"null",
			},
		},
		{
			name:    "JSON with boolean",
			content: `{"flag":false}`,
			contains: []string{
				"flag",
				"false",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.Format(tt.content)
			
			if got == "" {
				t.Error("Format() returned empty string")
			}

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Format() missing %q in output", want)
				}
			}
		})
	}
}

func TestJSONFormatterInvalidJSON(t *testing.T) {
	formatter := NewJSONFormatter()

	invalidJSON := `{"invalid": json}`
	result := formatter.Format(invalidJSON)

	if result != invalidJSON {
		t.Error("Format() should return original content for invalid JSON")
	}
}

func TestDetectAndFormat(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		contentType string
		wantFormat  bool
	}{
		{
			name:        "JSON content",
			content:     `{"key":"value"}`,
			contentType: "application/json",
			wantFormat:  true,
		},
		{
			name:        "Plain text",
			content:     "Hello world",
			contentType: "text/plain",
			wantFormat:  false,
		},
		{
			name:        "Auto-detect JSON",
			content:     `{"auto":"detect"}`,
			contentType: "",
			wantFormat:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectAndFormat(tt.content, tt.contentType)

			if tt.wantFormat {
				if result == tt.content {
					t.Error("DetectAndFormat() should format content")
				}
			} else {
				if result != tt.content {
					t.Error("DetectAndFormat() should not format content")
				}
			}
		})
	}
}
