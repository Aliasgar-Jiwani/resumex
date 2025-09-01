package plugins

import (
	"strings"

	"github.com/Aliasgar-Jiwani/resumex/pkg/session"
)

type ResumeFunc func(sess *session.Session) string

var plugins = map[string]ResumeFunc{
	"wget":  wgetResume,
	"rsync": rsyncResume,
	"tar":   tarResume,
	"curl":  curlResume,
}

func GetResumeCommand(sess *session.Session) string {
	// Parse the original command to get the base command
	parts := strings.Fields(sess.Command)
	if len(parts) == 0 {
		return sess.Command
	}

	baseCmd := parts[0]
	
	// Remove path if present (e.g., /usr/bin/wget -> wget)
	if lastSlash := strings.LastIndex(baseCmd, "/"); lastSlash != -1 {
		baseCmd = baseCmd[lastSlash+1:]
	}

	// Check if we have a plugin for this command
	if resumeFunc, exists := plugins[baseCmd]; exists {
		return resumeFunc(sess)
	}

	// Default: just re-run the original command
	return sess.Command
}

func wgetResume(sess *session.Session) string {
	// Add continuation flag if not already present
	if !strings.Contains(sess.Command, "-c") && !strings.Contains(sess.Command, "--continue") {
		return sess.Command + " -c"
	}
	return sess.Command
}

func rsyncResume(sess *session.Session) string {
	// Add partial flag if not already present
	if !strings.Contains(sess.Command, "--partial") {
		return sess.Command + " --partial"
	}
	return sess.Command
}

func tarResume(sess *session.Session) string {
	// For tar, we'll check if it's an extraction or creation
	if strings.Contains(sess.Command, "-x") {
		// Extraction - add keep-newer-files flag
		if !strings.Contains(sess.Command, "--keep-newer-files") {
			return sess.Command + " --keep-newer-files"
		}
	} else if strings.Contains(sess.Command, "-c") {
		// Creation - this is more complex, for now just re-run
		// In a real implementation, you might want to check what's already archived
		return sess.Command
	}
	return sess.Command
}

func curlResume(sess *session.Session) string {
	// Add continuation flag if not already present
	if !strings.Contains(sess.Command, "-C") && !strings.Contains(sess.Command, "--continue-at") {
		return sess.Command + " -C -"
	}
	return sess.Command
}

// RegisterPlugin allows users to add custom resume logic for specific commands
func RegisterPlugin(command string, resumeFunc ResumeFunc) {
	plugins[command] = resumeFunc
}

// ListPlugins returns all registered plugins
func ListPlugins() map[string]bool {
	result := make(map[string]bool)
	for cmd := range plugins {
		result[cmd] = true
	}
	return result
}
