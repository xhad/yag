package tests

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
)

// TestBranchCommand tests branch creation and listing
func TestBranchCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_branch_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	defer func() {
		os.RemoveAll(tempDir)
	}()

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer func() {
		os.Chdir(originalDir)
	}()

	// Initialize repository structure
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

	// Create a test file
	testFile := filepath.Join(tempDir, "test_file.txt")
	err = os.WriteFile(testFile, []byte("Hello, YAG!"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a blob from the test file
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	// Store the blob
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

	// Update master branch to point to the new commit
	masterRefPath := filepath.Join(headsDir, storage.DefaultBranch)
	if err := os.WriteFile(masterRefPath, []byte(commit.ID()), 0644); err != nil {
		t.Fatalf("Failed to update master ref: %v", err)
	}

	// Update HEAD to point to master branch properly
	headContent = "ref: refs/heads/" + storage.DefaultBranch
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		t.Fatalf("Failed to update HEAD file: %v", err)
	}

	// Test creating a new branch
	newBranch := "test-branch"
	err = commands.BranchCommand([]string{newBranch})
	if err != nil {
		t.Fatalf("BranchCommand failed to create branch: %v", err)
	}

	// Verify the branch file was created
	branchPath := filepath.Join(tempDir, storage.YAGDir, storage.RefsDir, storage.HeadsDir, newBranch)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		t.Errorf("Branch file was not created at %s", branchPath)
	}

	// Test listing branches (should include master and test-branch)
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = commands.BranchCommand([]string{})
	if err != nil {
		t.Fatalf("BranchCommand failed to list branches: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read captured output: %v", err)
	}
	output := buf.String()

	if !strings.Contains(output, storage.DefaultBranch) || !strings.Contains(output, newBranch) {
		t.Errorf("Branch listing should contain '%s' and '%s', got: %s", storage.DefaultBranch, newBranch, output)
	}

	// Test error case: creating a branch without any commits
	subDir := filepath.Join(tempDir, "no-commits")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("Failed to change to subdirectory: %v", err)
	}

	err = commands.InitCommand([]string{})
	if err != nil {
		t.Fatalf("Failed to initialize repository in subdirectory: %v", err)
	}

	// Try to create a branch
	err = commands.BranchCommand([]string{"error-branch"})
	if err == nil {
		t.Errorf("BranchCommand should fail in a repository with no commits")
	}
}
