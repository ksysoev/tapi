package cmd

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ksysoev/tapi/pkg/openapi"
	"github.com/ksysoev/tapi/pkg/tui"
)

func runExplore(ctx context.Context, filePath, url string) error {
	var spec *openapi.Spec
	var err error

	if filePath != "" {
		spec, err = openapi.LoadFromFile(filePath)
	} else {
		spec, err = openapi.LoadFromURL(url)
	}

	if err != nil {
		return fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	model := tui.NewModel(spec)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}
