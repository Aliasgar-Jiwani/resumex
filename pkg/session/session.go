package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusRunning     Status = "running"
	StatusCompleted   Status = "completed"
	StatusInterrupted Status = "interrupted"
)

type Session struct {
	ID         string    `json:"id"`
	Command    string    `json:"command"`
	WorkingDir string    `json:"working_dir"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time,omitempty"`
	Status     Status    `json:"status"`
	ExitCode   int       `json:"exit_code,omitempty"`
	LogFile    string    `json:"log_file"`
}

// Create a new session
func NewSession(command, workingDir string) *Session {
	id := uuid.New().String()

	configDir, err := getConfigDir()
	if err != nil {
		panic(fmt.Sprintf("Cannot get config directory: %v", err))
	}

	// Ensure sessions and logs directories exist
	logsDir := filepath.Join(configDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(fmt.Sprintf("Cannot create logs directory: %v", err))
	}

	sessionsDir := filepath.Join(configDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		panic(fmt.Sprintf("Cannot create sessions directory: %v", err))
	}

	logFile := filepath.Join(logsDir, id+".log")

	return &Session{
		ID:         id,
		Command:    command,
		WorkingDir: workingDir,
		StartTime:  time.Now(),
		Status:     StatusRunning,
		LogFile:    logFile,
	}
}

// Save session state to disk
func (s *Session) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	sessionsDir := filepath.Join(configDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return err
	}

	sessionFile := filepath.Join(sessionsDir, s.ID+".json")

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sessionFile, data, 0644)
}

// Load a session by ID
func LoadSession(sessionID string) (*Session, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	sessionFile := filepath.Join(configDir, "sessions", sessionID+".json")

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, err
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, err
	}

	return &sess, nil
}

// Delete session by ID
func DeleteSession(sessionID string) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	sessionFile := filepath.Join(configDir, "sessions", sessionID+".json")
	return os.Remove(sessionFile)
}

// Status update helpers
func (s *Session) MarkAsCompleted(exitCode int) {
	s.Status = StatusCompleted
	s.ExitCode = exitCode
	s.EndTime = time.Now()
}

func (s *Session) MarkAsInterrupted() {
	s.Status = StatusInterrupted
	s.EndTime = time.Now()
}

func (s *Session) MarkAsRunning() {
	s.Status = StatusRunning
}

// Get ~/.resumex path
func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".resumex"), nil
}
