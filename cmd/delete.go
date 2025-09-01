package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [session-id]",
	Short: "Delete a session",
	Long:  `Remove a session and its associated log files.`,
	Args:  cobra.ExactArgs(1),
	Run:   deleteCommand,
}

func deleteCommand(cmd *cobra.Command, args []string) {
	sessionID := args[0]

	// Load session to get log file path
	sess, err := session.LoadSession(sessionID)
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
	if err := session.DeleteSession(sessionID); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Session %s deleted successfully.\n", sessionID)
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}