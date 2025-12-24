package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Version string
	AppName string
}

func InitCommand(info BuildInfo) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     info.AppName,
		Short:   "Terminal API Explorer - Interactive OpenAPI specification browser",
		Long:    `TAPI is a beautiful terminal-based OpenAPI specification explorer and API testing tool.`,
		Version: info.Version,
	}

	rootCmd.AddCommand(newExploreCommand())
	rootCmd.AddCommand(newValidateCommand())

	return rootCmd
}

func newExploreCommand() *cobra.Command {
	var (
		filePath string
		url      string
	)

	cmd := &cobra.Command{
		Use:   "explore",
		Short: "Explore OpenAPI specification in interactive TUI",
		Long:  `Launch an interactive terminal UI to browse and test API endpoints defined in an OpenAPI specification.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" && url == "" {
				return fmt.Errorf("either --file or --url must be specified")
			}
			if filePath != "" && url != "" {
				return fmt.Errorf("only one of --file or --url can be specified")
			}

			return runExplore(cmd.Context(), filePath, url)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to local OpenAPI specification file")
	cmd.Flags().StringVarP(&url, "url", "u", "", "URL to remote OpenAPI specification")

	return cmd
}

func newValidateCommand() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate OpenAPI specification",
		Long:  `Validate an OpenAPI specification file for correctness.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath == "" {
				return fmt.Errorf("--file must be specified")
			}
			return runValidate(filePath)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to OpenAPI specification file (required)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
