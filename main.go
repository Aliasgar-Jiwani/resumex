package main

import (
	"fmt"
	"os"

	"github.com/yourusername/resumex/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// File: go.mod
module github.com/yourusername/resumex

go 1.21

require (
	github.com/google/uuid v1.3.0
	github.com/spf13/cobra v1.7.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

// File: cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "resumex",
	Short: "A universal resumable command wrapper",
	Long: `resumex allows you to run any command with state management,
so you can resume interrupted commands intelligently.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Ensure config directory exists
	configDir, err := getConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
		os.Exit(1)
	}

	sessionsDir := fmt.Sprintf("%s/sessions", configDir)
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating sessions directory: %v\n", err)
		os.Exit(1)
	}
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.resumex", home), nil
}
