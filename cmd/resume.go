package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	inputID := args[0]

	// Resolve full session ID if short ID provided
	_, fullID, err := findSessionFile(inputID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load session
	sess, err := session.LoadSession(fullID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
		os.Exit(1)
	}

	if sess.Status == session.StatusCompleted {
		fmt.Printf("Session %s is already completed\n", fullID)
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
	resumeCmdStr := plugins.GetResumeCommand(sess)
	fmt.Printf("Resume command: %s\n", resumeCmdStr)

	// Parse the resume command
	parts := parseCommand(resumeCmdStr)
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

// findSessionFile resolves a short or full ID to the actual session JSON file
func findSessionFile(inputID string) (string, string, error) {
	configDir, err := session.GetConfigDir()
	if err != nil {
		return "", "", fmt.Errorf("cannot get config dir: %v", err)
	}
	sessionsDir := filepath.Join(configDir, "sessions")

	files, err := os.ReadDir(sessionsDir)
	if err != nil {
		return "", "", fmt.Errorf("cannot read sessions dir: %v", err)
	}

	matches := []string{}
	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, inputID) && strings.HasSuffix(name, ".json") {
			matches = append(matches, name)
		}
	}

	if len(matches) == 0 {
		return "", "", fmt.Errorf("no session found matching ID '%s'", inputID)
	} else if len(matches) > 1 {
		return "", "", fmt.Errorf("multiple sessions match ID '%s', use full ID", inputID)
	}

	fullFile := filepath.Join(sessionsDir, matches[0])
	fullID := strings.TrimSuffix(matches[0], ".json")
	return fullFile, fullID, nil
}

// Simple parser for commands
func parseCommand(cmdStr string) []string {
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
