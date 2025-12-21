package openapi

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestLoadFromURL(t *testing.T) {
	validSpec := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/test": {
				"get": {
					"summary": "Test endpoint",
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		}
	}`

	tests := []struct {
		name       string
		serverFunc func() *httptest.Server
		wantErr    bool
	}{
		{
			name: "valid JSON spec",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(validSpec))
				}))
			},
			wantErr: false,
		},
		{
			name: "server error",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			defer server.Close()

			spec, err := LoadFromURL(server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if spec == nil {
					t.Error("LoadFromURL() returned nil spec without error")
				}
				if spec.Title != "Test API" {
					t.Errorf("Expected title 'Test API', got '%s'", spec.Title)
				}
			}
		})
	}
}

func TestLoadFromURLInvalidURL(t *testing.T) {
	_, err := LoadFromURL("http://invalid-url-does-not-exist-12345.com/spec.json")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestConvertSpec(t *testing.T) {
	data, err := os.ReadFile("../../example-petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	spec, err := parseSpec(data)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	if spec.Title == "" {
		t.Error("Expected non-empty title")
	}

	if spec.Version == "" {
		t.Error("Expected non-empty version")
	}

	if len(spec.Servers) == 0 {
		t.Error("Expected at least one server")
	}

	if len(spec.Paths) == 0 {
		t.Error("Expected at least one path")
	}

	hasOperations := false
	for _, path := range spec.Paths {
		if len(path.Operations) > 0 {
			hasOperations = true
			break
		}
	}
	if !hasOperations {
		t.Error("Expected at least one operation")
	}
}

func TestConvertSpecOperations(t *testing.T) {
	data, err := os.ReadFile("../../example-petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	spec, err := parseSpec(data)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	methodsFound := make(map[string]bool)
	hasParameters := false
	hasResponses := false

	for _, path := range spec.Paths {
		for _, op := range path.Operations {
			methodsFound[op.Method] = true

			if len(op.Parameters) > 0 {
				hasParameters = true
				param := op.Parameters[0]
				if param.Name == "" {
					t.Error("Parameter has empty name")
				}
			}

			if len(op.Responses) > 0 {
				hasResponses = true
			}
		}
	}

	if !hasParameters {
		t.Error("Expected to find at least one parameter")
	}

	if !hasResponses {
		t.Error("Expected to find at least one response")
	}
}

func TestConvertSchemaWithNilRef(t *testing.T) {
	schema := convertSchema(nil)
	if schema != nil {
		t.Error("Expected nil schema for nil input")
	}
}

func TestParseSpecInvalidData(t *testing.T) {
	invalidData := []byte("this is not valid yaml or json")
	_, err := parseSpec(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data")
	}
}

func TestServerConversion(t *testing.T) {
	data, err := os.ReadFile("../../example-petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	spec, err := parseSpec(data)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	if len(spec.Servers) == 0 {
		t.Fatal("Expected at least one server")
	}

	server := spec.Servers[0]
	if server.URL == "" {
		t.Error("Expected non-empty server URL")
	}
}

func TestOperationTags(t *testing.T) {
	data, err := os.ReadFile("../../example-petstore.yaml")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	spec, err := parseSpec(data)
	if err != nil {
		t.Fatalf("Failed to parse spec: %v", err)
	}

	foundTags := false
	for _, path := range spec.Paths {
		for _, op := range path.Operations {
			if len(op.Tags) > 0 {
				foundTags = true
				break
			}
		}
	}

	if !foundTags {
		t.Log("Note: No tags found in operations (not an error, just informational)")
	}
}
