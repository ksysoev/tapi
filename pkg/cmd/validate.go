package cmd

import (
	"fmt"

	"github.com/ksysoev/tapi/pkg/openapi"
)

func runValidate(filePath string) error {
	spec, err := openapi.LoadFromFile(filePath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("✓ OpenAPI specification is valid\n")
	fmt.Printf("  Title: %s\n", spec.Title)
	fmt.Printf("  Version: %s\n", spec.Version)
	fmt.Printf("  Endpoints: %d\n", len(spec.Paths))

	return nil
}
