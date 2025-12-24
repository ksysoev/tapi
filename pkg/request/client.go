package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ResponseMsg struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       string
	Error      error
}

func Send(baseURL, path, method string, params map[string]string, body string) tea.Cmd {
	return func() tea.Msg {
		fullURL := buildURL(baseURL, path, params)

		var reqBody io.Reader
		if body != "" {
			reqBody = bytes.NewBufferString(body)
		}

		req, err := http.NewRequest(method, fullURL, reqBody)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("failed to create request: %w", err)}
		}

		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}
		req.Header.Set("Accept", "application/json")

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("request failed: %w", err)}
		}
		defer func() { _ = resp.Body.Close() }()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return ResponseMsg{Error: fmt.Errorf("failed to read response: %w", err)}
		}

		return ResponseMsg{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    resp.Header,
			Body:       string(respBody),
		}
	}
}

func buildURL(baseURL, path string, params map[string]string) string {
	fullPath := strings.TrimSuffix(baseURL, "/") + path

	// Handle path parameters
	pathParams := make(map[string]bool)
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		if strings.Contains(fullPath, placeholder) {
			fullPath = strings.ReplaceAll(fullPath, placeholder, value)
			pathParams[key] = true
		}
	}

	// Handle query parameters
	queryParams := url.Values{}
	for key, value := range params {
		if !pathParams[key] && value != "" {
			queryParams.Add(key, value)
		}
	}

	if len(queryParams) > 0 {
		fullPath += "?" + queryParams.Encode()
	}

	return fullPath
}
