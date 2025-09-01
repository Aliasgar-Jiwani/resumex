package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Aliasgar-Jiwani/resumex/pkg/executor"
	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

var runCmd = &cobra.Command{
	Use:   "run [command and args...]",
	Short: "Run a command with resumable session tracking",
	Long:  `Run a command and save its state for potential resumption if interrupted.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runCommand,
}

func runCommand(cmd *cobra.Command, args []string) {
	// Create a new session
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	sess := session.NewSession(strings.Join(args, " "), wd)

	// Save session metadata
	if err := sess.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting session: %s\n", sess.ID)
	fmt.Printf("Command: %s\n", sess.Command)

	// Execute the command
	exec := executor.New(sess)
	exitCode, err := exec.Run(args[0], args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		sess.MarkAsInterrupted()
		sess.Save()
		os.Exit(1)
	}

	// Update session status
	if exitCode == 0 {
		sess.MarkAsCompleted()
	} else {
		sess.MarkAsInterrupted()
	}
	sess.Save()

	fmt.Printf("\nSession %s finished with exit code: %d\n", sess.ID, exitCode)
	os.Exit(exitCode)
}

func init() {
	rootCmd.AddCommand(runCmd)
}