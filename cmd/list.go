package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved sessions",
	Long:  `Display all saved sessions with their status and metadata.`,
	Run:   listCommand,
}

func listCommand(cmd *cobra.Command, args []string) {
	configDir, err := getConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config directory: %v\n", err)
		os.Exit(1)
	}

	sessionsDir := fmt.Sprintf("%s/sessions", configDir)
	
	// Read all session files
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", sessionsDir))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading sessions: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No sessions found.")
		return
	}

	// Load all sessions
	sessions := make([]*session.Session, 0, len(files))
	for _, file := range files {
		sessionID := filepath.Base(file)
		sessionID = sessionID[:len(sessionID)-5] // Remove .json extension
		
		sess, err := session.LoadSession(sessionID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not load session %s: %v\n", sessionID, err)
			continue
		}
		sessions = append(sessions, sess)
	}

	// Sort by start time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.After(sessions[j].StartTime)
	})

	// Display sessions in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SESSION ID\tSTATUS\tCOMMAND\tSTART TIME\tWORKING DIR")
	fmt.Fprintln(w, "----------\t------\t-------\t----------\t-----------")

	for _, sess := range sessions {
		// Truncate command if too long
		command := sess.Command
		if len(command) > 50 {
			command = command[:47] + "..."
		}

		// Format start time
		startTime := sess.StartTime.Format("2006-01-02 15:04")

		// Truncate working directory if too long
		workingDir := sess.WorkingDir
		if len(workingDir) > 30 {
			workingDir = "..." + workingDir[len(workingDir)-27:]
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			sess.ID[:8], // Show first 8 chars of ID
			sess.Status,
			command,
			startTime,
			workingDir)
	}

	w.Flush()
}

func init() {
	rootCmd.AddCommand(listCmd)
}