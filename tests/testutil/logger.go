package testutil

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Logger provides pretty logging for tests
type Logger struct {
	t *testing.T
}

// NewLogger creates a new logger for the given test
func NewLogger(t *testing.T) *Logger {
	return &Logger{t: t}
}

// StartTest logs the beginning of a test
func (l *Logger) StartTest() {
	testName := l.t.Name()
	l.t.Logf("\n%s=== STARTING TEST: %s ===%s\n", colorBlue, testName, colorReset)
	l.t.Logf("%s%s%s\n", colorBlue, strings.Repeat("=", len(testName)+18), colorReset)
}

// EndTest logs the end of a test
func (l *Logger) EndTest() {
	testName := l.t.Name()
	l.t.Logf("\n%s=== COMPLETED TEST: %s ===%s\n", colorGreen, testName, colorReset)
	l.t.Logf("%s%s%s\n", colorGreen, strings.Repeat("=", len(testName)+20), colorReset)
}

// Section logs a section header
func (l *Logger) Section(name string) {
	l.t.Logf("\n%s--- %s ---%s\n", colorCyan, name, colorReset)
}

// Info logs an informational message
func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.t.Logf("%s➤ %s%s", colorWhite, msg, colorReset)
}

// Success logs a success message
func (l *Logger) Success(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.t.Logf("%s✓ %s%s", colorGreen, msg, colorReset)
}

// Warning logs a warning message
func (l *Logger) Warning(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.t.Logf("%s⚠ %s%s", colorYellow, msg, colorReset)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	l.t.Logf("%s✗ %s%s", colorRed, msg, colorReset)
}

// Action logs an action being performed
func (l *Logger) Action(action, target string) {
	l.t.Logf("%s→ %s:%s %s", colorPurple, action, colorReset, target)
}

// Command logs a command being executed
func (l *Logger) Command(cmd string) {
	l.t.Logf("%s$ %s%s", colorYellow, cmd, colorReset)
}

// Timing logs the time taken for an operation
// @notice Logs the elapsed time for an operation with blue text formatting
// @param operation The name of the operation that was timed
// @param start The time.Time when the operation started
func (l *Logger) Timing(operation string, start time.Time) {
	elapsed := time.Since(start)
	l.t.Logf("%s%s:%s %s", colorBlue, operation, elapsed, colorReset)
}

// Separator prints a separator line
func (l *Logger) Separator() {
	l.t.Logf("\n%s%s%s\n", colorBlue, strings.Repeat("-", 80), colorReset)
}

// File logs information about a file
func (l *Logger) File(path, action string) {
	l.t.Logf("%s%s:%s %s", colorCyan, action, colorReset, path)
}

// Repository logs information about a repository operation
func (l *Logger) Repository(operation string, details string) {
	l.t.Logf("%s%s:%s %s", colorGreen, operation, colorReset, details)
}
