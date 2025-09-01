package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Aliasgar-Jiwani/resumex/pkg/executor"
	"github.com/Aliasgar-Jiwani/resumex/pkg/plugins"
	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

var resumeCmd = &cobra.Command{
	Use:   "resume [session-id]",
	Short: "Resume a previously interrupted session",
	Long:  `Resume a command that was previously interrupted using the session ID.`,
	Args:  cobra.ExactArgs(1),
	Run:   resumeCommand,
}

func resumeCommand(cmd *cobra.Command, args []string) {
	sessionID := args[0]

	// Load session
	sess, err := session.LoadSession(sessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
		os.Exit(1)
	}

	if sess.Status == session.StatusCompleted {
		fmt.Printf("Session %s is already completed\n", sessionID)
		return
	}

	fmt.Printf("Resuming session: %s\n", sess.ID)
	fmt.Printf("Original command: %s\n", sess.Command)

	// Change to original working directory
	if err := os.Chdir(sess.WorkingDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error changing to working directory %s: %v\n", sess.WorkingDir, err)
		os.Exit(1)
	}

	// Get resume command using plugins
	resumeCmd := plugins.GetResumeCommand(sess)
	fmt.Printf("Resume command: %s\n", resumeCmd)

	// Parse the resume command
	parts := parseCommand(resumeCmd)
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "Invalid resume command\n")
		os.Exit(1)
	}

	// Execute the resume command
	exec := executor.New(sess)
	exitCode, err := exec.Run(parts[0], parts[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing resume command: %v\n", err)
		sess.MarkAsInterrupted()
		sess.Save()
		os.Exit(1)
	}

	// Update session status
	if exitCode == 0 {
		sess.MarkAsCompleted(exitCode)
	} else {
		sess.MarkAsInterrupted()
	}
	sess.Save()

	fmt.Printf("\nSession %s resumed and finished with exit code: %d\n", sess.ID, exitCode)
	os.Exit(exitCode)
}

func parseCommand(cmdStr string) []string {
	// Simple command parsing - in production you might want to use a proper shell parser
	parts := []string{}
	current := ""
	inQuotes := false

	for i, char := range cmdStr {
		if char == '"' || char == '\'' {
			inQuotes = !inQuotes
		} else if char == ' ' && !inQuotes {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}

		if i == len(cmdStr)-1 && current != "" {
			parts = append(parts, current)
		}
	}

	return parts
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
