package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time" 

	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

type Executor struct {
	session *session.Session
	logFile *os.File
}

func New(sess *session.Session) *Executor {
	return &Executor{
		session: sess,
	}
}

func (e *Executor) Run(command string, args ...string) (int, error) {
	// Open log file
	logFile, err := os.OpenFile(e.session.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 1, fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()
	e.logFile = logFile

	// Log session start
	fmt.Fprintf(logFile, "=== Session %s started at %s ===\n", e.session.ID, e.session.StartTime.Format(time.RFC3339))
	fmt.Fprintf(logFile, "Command: %s %v\n", command, args)
	fmt.Fprintf(logFile, "Working Directory: %s\n", e.session.WorkingDir)
	fmt.Fprintf(logFile, "========================================\n")

	// Create command
	cmd := exec.Command(command, args...)
	cmd.Dir = e.session.WorkingDir

	// Set up pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return 1, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 1, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return 1, fmt.Errorf("failed to start command: %w", err)
	}

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to receive command completion
	done := make(chan error, 1)

	// Start goroutines to handle output
	go e.handleOutput(stdoutPipe, "STDOUT")
	go e.handleOutput(stderrPipe, "STDERR")

	// Wait for command completion in a goroutine
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for either completion or signal
	select {
	case err := <-done:
		// Command completed
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				return exitError.ExitCode(), nil
			}
			return 1, err
		}
		return 0, nil
	case sig := <-sigChan:
		// Signal received, terminate the command
		fmt.Printf("\nReceived signal %s, terminating command...\n", sig)
		
		// Send signal to the process group
		if cmd.Process != nil {
			cmd.Process.Signal(sig)
		}

		// Wait a bit for graceful shutdown
		select {
		case <-done:
			// Command finished gracefully
		case <-time.After(5 * time.Second):
			// Force kill if not finished
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}

		return 130, fmt.Errorf("command interrupted by signal %s", sig)
	}
}

func (e *Executor) handleOutput(pipe io.ReadCloser, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Write to log file
		fmt.Fprintf(e.logFile, "[%s] %s\n", prefix, line)
		
		// Write to stdout
		fmt.Println(line)
	}
}