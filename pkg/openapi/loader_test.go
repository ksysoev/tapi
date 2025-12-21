package openapi

import (
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid yaml spec",
			path:    "../../example-petstore.yaml",
			wantErr: false,
		},
		{
			name:    "non-existent file",
			path:    "non-existent.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, err := LoadFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && spec == nil {
				t.Error("LoadFromFile() returned nil spec without error")
			}
			if !tt.wantErr {
				if spec.Title == "" {
					t.Error("LoadFromFile() spec has empty title")
				}
				if len(spec.Paths) == 0 {
					t.Error("LoadFromFile() spec has no paths")
				}
			}
		})
	}
}
