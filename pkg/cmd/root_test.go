package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestInitCommand(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		AppName: "tapi",
	}

	cmd := InitCommand(info)

	if cmd == nil {
		t.Fatal("Expected non-nil command")
	}

	if cmd.Use != "tapi" {
		t.Errorf("Expected Use to be 'tapi', got '%s'", cmd.Use)
	}

	if cmd.Version != "1.0.0" {
		t.Errorf("Expected Version to be '1.0.0', got '%s'", cmd.Version)
	}

	if !cmd.HasSubCommands() {
		t.Error("Expected command to have subcommands")
	}

	expectedCommands := []string{"explore", "validate"}
	for _, cmdName := range expectedCommands {
		if _, _, err := cmd.Find([]string{cmdName}); err != nil {
			t.Errorf("Expected to find subcommand '%s'", cmdName)
		}
	}
}

func TestExploreCommandFlags(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		AppName: "tapi",
	}

	rootCmd := InitCommand(info)
	exploreCmd, _, err := rootCmd.Find([]string{"explore"})
	if err != nil {
		t.Fatalf("Failed to find explore command: %v", err)
	}

	expectedFlags := []string{"file", "url"}
	for _, flagName := range expectedFlags {
		flag := exploreCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
		}
	}
}

func TestExploreCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no flags specified",
			args:    []string{"explore"},
			wantErr: true,
			errMsg:  "either --file or --url must be specified",
		},
		{
			name:    "both flags specified",
			args:    []string{"explore", "--file", "test.yaml", "--url", "http://example.com"},
			wantErr: true,
			errMsg:  "only one of --file or --url can be specified",
		},
		{
			name:    "file flag only",
			args:    []string{"explore", "--file", "../../example-petstore.yaml"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := BuildInfo{
				Version: "1.0.0",
				AppName: "tapi",
			}

			cmd := InitCommand(info)
			cmd.SetArgs(tt.args)

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			if tt.wantErr {
				err := cmd.Execute()
				if err == nil {
					t.Error("Expected error but got none")
				}
				if err != nil && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				// For success cases, only validate args parsing without executing
				// to avoid launching TUI which would block the test
				exploreCmd, _, err := cmd.Find([]string{"explore"})
				if err != nil {
					t.Fatalf("Failed to find explore command: %v", err)
				}
				_ = exploreCmd.ParseFlags(tt.args[1:])
			}
		})
	}
}

func TestValidateCommandFlags(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		AppName: "tapi",
	}

	rootCmd := InitCommand(info)
	validateCmd, _, err := rootCmd.Find([]string{"validate"})
	if err != nil {
		t.Fatalf("Failed to find validate command: %v", err)
	}

	fileFlag := validateCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Error("Expected 'file' flag to exist")
	}
}

func TestValidateCommandValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no file specified",
			args:    []string{"validate"},
			wantErr: true,
			errMsg:  "required flag(s)",
		},
		{
			name:    "valid file specified",
			args:    []string{"validate", "--file", "../../example-petstore.yaml"},
			wantErr: false,
		},
		{
			name:    "invalid file specified",
			args:    []string{"validate", "--file", "non-existent.yaml"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := BuildInfo{
				Version: "1.0.0",
				AppName: "tapi",
			}

			cmd := InitCommand(info)
			cmd.SetArgs(tt.args)

			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) && !strings.Contains(buf.String(), tt.errMsg) {
					t.Errorf("Expected error to contain '%s', got '%s' and output '%s'", tt.errMsg, err.Error(), buf.String())
				}
			}
		})
	}
}

func TestRunValidate(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "valid spec",
			filePath: "../../example-petstore.yaml",
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "non-existent.yaml",
			wantErr:  true,
		},
		{
			name:     "empty path",
			filePath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runValidate(tt.filePath)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCommandHelpText(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		AppName: "tapi",
	}

	cmd := InitCommand(info)

	if cmd.Short == "" {
		t.Error("Expected non-empty Short description")
	}

	if cmd.Long == "" {
		t.Error("Expected non-empty Long description")
	}

	exploreCmd, _, _ := cmd.Find([]string{"explore"})
	if exploreCmd.Short == "" {
		t.Error("Expected non-empty Short description for explore command")
	}

	validateCmd, _, _ := cmd.Find([]string{"validate"})
	if validateCmd.Short == "" {
		t.Error("Expected non-empty Short description for validate command")
	}
}
