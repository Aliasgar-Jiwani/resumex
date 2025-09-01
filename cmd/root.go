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
