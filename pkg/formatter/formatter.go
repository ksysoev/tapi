// Package formatter provides automatic format detection and syntax highlighting
// for API response bodies. It currently supports JSON formatting with syntax
// highlighting and is designed to be easily extensible for additional formats.
package formatter

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ContentType string

const (
	ContentTypeJSON    ContentType = "json"
	ContentTypeUnknown ContentType = "unknown"
)

type Formatter interface {
	Format(content string) string
	CanHandle(content string, contentType string) bool
}

var (
	jsonKeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	jsonStringStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	jsonNumberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("141"))
	jsonBoolStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	jsonNullStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type jsonFormatter struct{}

func NewJSONFormatter() Formatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) CanHandle(content string, contentType string) bool {
	if strings.Contains(contentType, "json") {
		return true
	}

	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 {
		return false
	}

	firstChar := trimmed[0]
	return firstChar == '{' || firstChar == '['
}

func (f *jsonFormatter) Format(content string) string {
	var parsed interface{}
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return content
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return content
	}

	return f.highlight(string(formatted))
}

func (f *jsonFormatter) highlight(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		result.WriteString(f.highlightLine(line))
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

func (f *jsonFormatter) highlightLine(line string) string {
	trimmed := strings.TrimSpace(line)
	
	if trimmed == "{" || trimmed == "}" || trimmed == "[" || trimmed == "]" || 
	   trimmed == "{}" || trimmed == "[]" || trimmed == "," {
		return line
	}

	if idx := strings.Index(line, ":"); idx != -1 {
		indent := line[:len(line)-len(strings.TrimLeft(line, " "))]
		afterIndent := line[len(indent):]
		
		colonIdx := strings.Index(afterIndent, ":")
		if colonIdx == -1 {
			return line
		}

		key := strings.TrimSpace(afterIndent[:colonIdx])
		rest := afterIndent[colonIdx+1:]
		value := strings.TrimSpace(rest)

		hasComma := strings.HasSuffix(value, ",")
		if hasComma {
			value = strings.TrimSuffix(value, ",")
			value = strings.TrimSpace(value)
		}

		styledKey := jsonKeyStyle.Render(key)
		styledValue := f.styleValue(value)

		result := indent + styledKey + ": " + styledValue
		if hasComma {
			result += ","
		}

		return result
	}

	return line
}

func (f *jsonFormatter) styleValue(value string) string {
	if value == "" {
		return value
	}

	if value == "true" || value == "false" {
		return jsonBoolStyle.Render(value)
	}

	if value == "null" {
		return jsonNullStyle.Render(value)
	}

	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		return jsonStringStyle.Render(value)
	}

	if _, err := json.Number(value).Float64(); err == nil {
		return jsonNumberStyle.Render(value)
	}

	return value
}

func DetectAndFormat(content string, contentType string) string {
	formatters := []Formatter{
		NewJSONFormatter(),
	}

	for _, formatter := range formatters {
		if formatter.CanHandle(content, contentType) {
			return formatter.Format(content)
		}
	}

	return content
}
