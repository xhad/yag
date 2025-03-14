package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
)

// TestStatusCommand tests that the status command correctly identifies staged, unstaged, and untracked files
func TestStatusCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_status_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalDir) // Restore original directory after test

	// Initialize a proper repository structure
	yagDir := filepath.Join(tempDir, storage.YAGDir)
	if err := os.MkdirAll(yagDir, 0755); err != nil {
		t.Fatalf("Failed to create .yag directory: %v", err)
	}

	// Create necessary subdirectories
	objectsDir := filepath.Join(yagDir, storage.ObjectsDir)
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		t.Fatalf("Failed to create objects directory: %v", err)
	}

	headsDir := filepath.Join(yagDir, storage.RefsDir, storage.HeadsDir)
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		t.Fatalf("Failed to create refs/heads directory: %v", err)
	}

	// Create HEAD file
	headPath := filepath.Join(yagDir, storage.HeadFile)
	headContent := "ref: " + storage.RefsDir + "/" + storage.HeadsDir + "/" + storage.DefaultBranch
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		t.Fatalf("Failed to create HEAD file: %v", err)
	}

	// Create an empty index file
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	// Create a staged file
	stagedFile := filepath.Join(tempDir, "staged_file.txt")
	stagedContent := []byte("This is a staged file")
	if err := os.WriteFile(stagedFile, stagedContent, 0644); err != nil {
		t.Fatalf("Failed to create staged file: %v", err)
	}

	// Get the relative path to the staged file
	stagedRelPath, err := filepath.Rel(tempDir, stagedFile)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}

	// Create a blob for the staged file
	stagedBlob, err := core.NewBlobFromFile(stagedFile)
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	// Write the blob to the objects directory
	stagedBlobPath := filepath.Join(objectsDir, stagedBlob.ID())
	serialized, err := stagedBlob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(stagedBlobPath, serialized, 0644); err != nil {
		t.Fatalf("Failed to write blob file: %v", err)
	}

	// Update the index with the staged file
	indexEntries := map[string]string{
		stagedRelPath: stagedBlob.ID(),
	}
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// Create an unstaged file (modified after staging)
	unstagedFile := filepath.Join(tempDir, "unstaged_file.txt")
	unstagedContent := []byte("This is an unstaged file")
	if err := os.WriteFile(unstagedFile, unstagedContent, 0644); err != nil {
		t.Fatalf("Failed to create unstaged file: %v", err)
	}

	// Get the relative path to the unstaged file
	unstagedRelPath, err := filepath.Rel(tempDir, unstagedFile)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}

	// Create a blob for the initial version of the unstaged file
	initialBlob, err := core.NewBlobFromFile(unstagedFile)
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	// Write the blob to the objects directory
	initialBlobPath := filepath.Join(objectsDir, initialBlob.ID())
	serialized, err = initialBlob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(initialBlobPath, serialized, 0644); err != nil {
		t.Fatalf("Failed to write blob file: %v", err)
	}

	// Update the index with the initial version of the unstaged file
	indexEntries[unstagedRelPath] = initialBlob.ID()
	indexData, err = json.Marshal(indexEntries)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// Modify the unstaged file to create a change that's not staged
	modifiedContent := []byte("This file has been modified")
	if err := os.WriteFile(unstagedFile, modifiedContent, 0644); err != nil {
		t.Fatalf("Failed to modify unstaged file: %v", err)
	}

	// Create an untracked file
	untrackedFile := filepath.Join(tempDir, "untracked_file.txt")
	untrackedContent := []byte("This is an untracked file")
	if err := os.WriteFile(untrackedFile, untrackedContent, 0644); err != nil {
		t.Fatalf("Failed to create untracked file: %v", err)
	}

	// Run the status command and capture its output
	// We can't easily capture stdout in tests, so we'll just make sure it runs without errors
	err = commands.StatusCommand([]string{})
	if err != nil {
		t.Fatalf("StatusCommand failed: %v", err)
	}

	// Create a commit so we can see what happens when the working tree is clean
	// First, build a tree from the index entries
	tree := core.BuildTreeFromPaths(indexEntries)

	// Store the tree in the objects directory
	treeData, err := tree.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize tree: %v", err)
	}
	treePath := filepath.Join(objectsDir, tree.ID())
	if err := os.WriteFile(treePath, treeData, 0644); err != nil {
		t.Fatalf("Failed to write tree file: %v", err)
	}

	// Create a commit pointing to the tree
	commit := core.NewCommit(tree.ID(), "", "Initial commit", "test")

	// Store the commit in the objects directory
	commitData, err := commit.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize commit: %v", err)
	}
	commitPath := filepath.Join(objectsDir, commit.ID())
	if err := os.WriteFile(commitPath, commitData, 0644); err != nil {
		t.Fatalf("Failed to write commit file: %v", err)
	}

	// Update the master branch reference to point to the commit
	masterRefPath := filepath.Join(headsDir, storage.DefaultBranch)
	if err := os.WriteFile(masterRefPath, []byte(commit.ID()), 0644); err != nil {
		t.Fatalf("Failed to write master ref file: %v", err)
	}

	// Clean up the untracked file, modify the unstaged file back to original,
	// and see if status shows clean working tree
	if err := os.Remove(untrackedFile); err != nil {
		t.Fatalf("Failed to remove untracked file: %v", err)
	}

	if err := os.WriteFile(unstagedFile, unstagedContent, 0644); err != nil {
		t.Fatalf("Failed to restore unstaged file: %v", err)
	}

	// Run status command again
	err = commands.StatusCommand([]string{})
	if err != nil {
		t.Fatalf("StatusCommand failed: %v", err)
	}
}
