package request

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		params   map[string]string
		expected string
	}{
		{
			name:     "simple path",
			baseURL:  "https://api.example.com",
			path:     "/users",
			params:   map[string]string{},
			expected: "https://api.example.com/users",
		},
		{
			name:     "path parameter",
			baseURL:  "https://api.example.com",
			path:     "/users/{id}",
			params:   map[string]string{"id": "123"},
			expected: "https://api.example.com/users/123",
		},
		{
			name:     "multiple path parameters",
			baseURL:  "https://api.example.com",
			path:     "/users/{userId}/posts/{postId}",
			params:   map[string]string{"userId": "123", "postId": "456"},
			expected: "https://api.example.com/users/123/posts/456",
		},
		{
			name:     "query parameters",
			baseURL:  "https://api.example.com",
			path:     "/users",
			params:   map[string]string{"page": "1", "limit": "10"},
			expected: "https://api.example.com/users?limit=10&page=1",
		},
		{
			name:     "mixed path and query parameters",
			baseURL:  "https://api.example.com",
			path:     "/users/{id}",
			params:   map[string]string{"id": "123", "include": "posts"},
			expected: "https://api.example.com/users/123?include=posts",
		},
		{
			name:     "base URL with trailing slash",
			baseURL:  "https://api.example.com/",
			path:     "/users",
			params:   map[string]string{},
			expected: "https://api.example.com/users",
		},
		{
			name:     "empty query parameter",
			baseURL:  "https://api.example.com",
			path:     "/users",
			params:   map[string]string{"page": "1", "empty": ""},
			expected: "https://api.example.com/users?page=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildURL(tt.baseURL, tt.path, tt.params)
			if got != tt.expected {
				t.Errorf("buildURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		params         map[string]string
		body           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
	}{
		{
			name:   "successful GET request",
			method: "GET",
			path:   "/users",
			params: map[string]string{},
			body:   "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET request, got %s", r.Method)
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"message": "success"})
			},
			wantErr: false,
		},
		{
			name:   "successful POST request with body",
			method: "POST",
			path:   "/users",
			params: map[string]string{},
			body:   `{"name":"John"}`,
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Expected POST request, got %s", r.Method)
				}
				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
			},
			wantErr: false,
		},
		{
			name:   "request with path parameter",
			method: "GET",
			path:   "/users/{id}",
			params: map[string]string{"id": "123"},
			body:   "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/users/123" {
					t.Errorf("Expected path /users/123, got %s", r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"id": "123"})
			},
			wantErr: false,
		},
		{
			name:   "request with query parameters",
			method: "GET",
			path:   "/users",
			params: map[string]string{"page": "1", "limit": "10"},
			body:   "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				page := r.URL.Query().Get("page")
				limit := r.URL.Query().Get("limit")
				if page != "1" || limit != "10" {
					t.Errorf("Expected page=1&limit=10, got page=%s&limit=%s", page, limit)
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"count": "10"})
			},
			wantErr: false,
		},
		{
			name:   "server error response",
			method: "GET",
			path:   "/error",
			params: map[string]string{},
			body:   "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
			},
			wantErr: false,
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			path:   "/users/{id}",
			params: map[string]string{"id": "123"},
			body:   "",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Expected DELETE request, got %s", r.Method)
				}
				w.WriteHeader(http.StatusNoContent)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			cmd := Send(server.URL, tt.path, tt.method, tt.params, tt.body)
			msg := cmd()

			responseMsg, ok := msg.(ResponseMsg)
			if !ok {
				t.Fatal("Expected ResponseMsg type")
			}

			if tt.wantErr && responseMsg.Error == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && responseMsg.Error != nil {
				t.Errorf("Unexpected error: %v", responseMsg.Error)
			}
		})
	}
}

func TestSendInvalidURL(t *testing.T) {
	cmd := Send("http://invalid-url-that-does-not-exist-12345.com", "/test", "GET", nil, "")
	msg := cmd()

	responseMsg, ok := msg.(ResponseMsg)
	if !ok {
		t.Fatal("Expected ResponseMsg type")
	}

	if responseMsg.Error == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestSendReturnsTeaCmd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cmd := Send(server.URL, "/test", "GET", nil, "")

	if cmd == nil {
		t.Fatal("Expected non-nil tea.Cmd")
	}

	var _ = cmd
}

func TestSendHeadersSet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("Expected Accept header to be application/json, got %s", accept)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cmd := Send(server.URL, "/test", "GET", nil, "")
	cmd()
}
