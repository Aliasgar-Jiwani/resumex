package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [session-id]",
	Short: "Delete a session",
	Long:  `Remove a session and its associated log files.`,
	Args:  cobra.ExactArgs(1),
	Run:   deleteCommand,
}

func deleteCommand(cmd *cobra.Command, args []string) {
	inputID := args[0]

	// Find the full session ID (supports short IDs)
	fullID, err := findFullSessionID(inputID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding session: %v\n", err)
		os.Exit(1)
	}

	// Load session to get log file path
	sess, err := session.LoadSession(fullID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
		os.Exit(1)
	}

	// Delete log file if it exists
	if sess.LogFile != "" {
		if err := os.Remove(sess.LogFile); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Warning: Could not delete log file %s: %v\n", sess.LogFile, err)
		}
	}

	// Delete session file
	if err := session.DeleteSession(fullID); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Session %s deleted successfully.\n", inputID)
}

// findFullSessionID finds the full UUID filename for a given short or full ID
func findFullSessionID(shortID string) (string, error) {
	configDir, err := session.GetConfigDir()
	if err != nil {
		return "", err
	}

	sessionsDir := filepath.Join(configDir, "sessions")
	files, err := os.ReadDir(sessionsDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, shortID) {
			return strings.TrimSuffix(name, ".json"), nil
		}
	}

	return "", fmt.Errorf("session with ID starting %s not found", shortID)
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
