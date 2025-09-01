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

func NewSession(command, workingDir string) *Session {
	id := uuid.New().String()
	
	configDir, err := getConfigDir()
	if err != nil {
		panic(fmt.Sprintf("Cannot get config directory: %v", err))
	}

	logFile := fmt.Sprintf("%s/logs/%s.log", configDir, id)
	
	// Ensure logs directory exists
	logsDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(fmt.Sprintf("Cannot create logs directory: %v", err))
	}

	return &Session{
		ID:         id,
		Command:    command,
		WorkingDir: workingDir,
		StartTime:  time.Now(),
		Status:     StatusRunning,
		LogFile:    logFile,
	}
}

func (s *Session) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	sessionsDir := fmt.Sprintf("%s/sessions", configDir)
	sessionFile := fmt.Sprintf("%s/%s.json", sessionsDir, s.ID)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sessionFile, data, 0644)
}

func LoadSession(sessionID string) (*Session, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	sessionFile := fmt.Sprintf("%s/sessions/%s.json", configDir, sessionID)

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

func DeleteSession(sessionID string) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	sessionFile := fmt.Sprintf("%s/sessions/%s.json", configDir, sessionID)
	return os.Remove(sessionFile)
}

func (s *Session) MarkAsCompleted() {
	s.Status = StatusCompleted
	s.ExitCode = 0
	s.EndTime = time.Now()
}

func (s *Session) MarkAsInterrupted() {
	s.Status = StatusInterrupted
	s.EndTime = time.Now()
}

func (s *Session) MarkAsRunning() {
	s.Status = StatusRunning
}

func getConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/.resumex", home), nil
}