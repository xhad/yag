package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/storage"
)

// TestInitCommand tests the initialization of a new repository
func TestInitCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_init_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Test initialization
	args := []string{tempDir}
	err = commands.InitCommand(args)
	if err != nil {
		t.Fatalf("InitCommand failed: %v", err)
	}

	// Verify the .yag directory was created
	yagDir := filepath.Join(tempDir, storage.YAGDir)
	if _, err := os.Stat(yagDir); os.IsNotExist(err) {
		t.Errorf(".yag directory was not created at %s", yagDir)
	}

	// Verify objects directory exists
	objectsDir := filepath.Join(yagDir, "objects")
	if _, err := os.Stat(objectsDir); os.IsNotExist(err) {
		t.Errorf("objects directory was not created at %s", objectsDir)
	}

	// Verify refs directory exists
	refsDir := filepath.Join(yagDir, "refs")
	if _, err := os.Stat(refsDir); os.IsNotExist(err) {
		t.Errorf("refs directory was not created at %s", refsDir)
	}

	// Test initializing an existing repository (should not error)
	err = commands.InitCommand(args)
	if err != nil {
		t.Errorf("Reinitializing repository should not error, got: %v", err)
	}
}
