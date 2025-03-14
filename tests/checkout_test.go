package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
)

// TestCheckoutCommand tests switching between branches
func TestCheckoutCommand(t *testing.T) {
	// Use the utility function to set up a test repository with a committed file
	tempDir, testFile, err := SetupTestRepo()
	if err != nil {
		t.Fatalf("Failed to set up test repository: %v", err)
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

	// Since the default GetIndexEntries implementation is just a placeholder,
	// we need to manually commit to have a valid commit to create a branch from

	// Create a blob from the test file
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	// Store the blob
	objectsDir := filepath.Join(tempDir, storage.YAGDir, "objects")
	objectPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(objectPath, serialized, 0644); err != nil {
		t.Fatalf("Failed to write blob to file: %v", err)
	}

	// Create the tree with the relative path to the file
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}
	treeEntries := map[string]string{
		relPath: blob.ID(),
	}
	tree := core.BuildTreeFromPaths(treeEntries)

	// Store the tree
	treeObjectPath := filepath.Join(objectsDir, tree.ID())
	treeSerialized, err := tree.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize tree: %v", err)
	}
	if err := os.WriteFile(treeObjectPath, treeSerialized, 0644); err != nil {
		t.Fatalf("Failed to write tree to file: %v", err)
	}

	// Create a commit
	author := "test-user"
	commitMessage := "Initial commit"
	commit := core.NewCommit(tree.ID(), "", commitMessage, author)

	// Store the commit
	commitObjectPath := filepath.Join(objectsDir, commit.ID())
	commitSerialized, err := commit.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize commit: %v", err)
	}
	if err := os.WriteFile(commitObjectPath, commitSerialized, 0644); err != nil {
		t.Fatalf("Failed to write commit to file: %v", err)
	}

	// Create the refs/heads directory structure if it doesn't exist
	headsDir := filepath.Join(tempDir, storage.YAGDir, storage.RefsDir, storage.HeadsDir)
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		t.Fatalf("Failed to create refs/heads directory: %v", err)
	}

	// Update master branch to point to the new commit
	masterRefPath := filepath.Join(headsDir, storage.DefaultBranch)
	if err := os.WriteFile(masterRefPath, []byte(commit.ID()), 0644); err != nil {
		t.Fatalf("Failed to update master ref: %v", err)
	}

	// Update HEAD to point to master branch properly
	headPath := filepath.Join(tempDir, storage.YAGDir, storage.HeadFile)
	headContent := "ref: refs/heads/" + storage.DefaultBranch
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		t.Fatalf("Failed to update HEAD file: %v", err)
	}

	// Create a new branch
	newBranch := "test-branch"
	err = commands.BranchCommand([]string{newBranch})
	if err != nil {
		t.Fatalf("Failed to create branch: %v", err)
	}

	// Test checkout the new branch
	err = commands.CheckoutCommand(newBranch)
	if err != nil {
		t.Fatalf("CheckoutCommand failed: %v", err)
	}

	// Verify HEAD points to the new branch
	headContentBytes, err := os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("Failed to read HEAD file: %v", err)
	}
	headContentStr := string(headContentBytes)

	expectedHeadContent := "ref: refs/heads/" + newBranch
	if !strings.Contains(headContentStr, newBranch) {
		t.Errorf("HEAD should point to %s, got: %s", expectedHeadContent, headContentStr)
	}

	// Test checkout back to master
	err = commands.CheckoutCommand(storage.DefaultBranch)
	if err != nil {
		t.Fatalf("CheckoutCommand failed to switch to master: %v", err)
	}

	// Verify HEAD now points to master
	headContentBytes, err = os.ReadFile(headPath)
	if err != nil {
		t.Fatalf("Failed to read HEAD file: %v", err)
	}
	headContentStr = string(headContentBytes)

	if !strings.Contains(headContentStr, storage.DefaultBranch) {
		t.Errorf("HEAD should point to %s, got: %s", storage.DefaultBranch, headContentStr)
	}

	// Test error case: checkout non-existent branch
	err = commands.CheckoutCommand("non-existent-branch")
	if err == nil {
		t.Errorf("CheckoutCommand should fail with non-existent branch")
	}

	// Test error case: empty branch name
	err = commands.CheckoutCommand("")
	if err == nil {
		t.Errorf("CheckoutCommand should fail with empty branch name")
	}
}
