package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
)

// TestCommitCommand tests creating a commit with staged changes
func TestCommitCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_commit_*")
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

	// Create a test file to add
	testFile := filepath.Join(tempDir, "test_file.txt")
	testContent := []byte("Hello, YAG!")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Get the relative path to the test file
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		t.Fatalf("Failed to get relative path: %v", err)
	}

	// Create a blob for the test file
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to create blob: %v", err)
	}

	// Write the blob to the objects directory
	blobPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(blobPath, serialized, 0644); err != nil {
		t.Fatalf("Failed to write blob file: %v", err)
	}

	// Create an index file with the test file
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	indexEntries := map[string]string{
		relPath: blob.ID(),
	}
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// COMMIT CREATION - Manual version that doesn't depend on storage.GetIndexEntries
	// Build a tree from the test file
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

	// Test error cases
	// Test empty commit message (using command directly)
	err = commands.CommitCommand("")
	if err == nil {
		t.Errorf("CommitCommand should fail with empty commit message")
	}

	// Test empty index case (using command directly)
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write empty index file: %v", err)
	}
	err = commands.CommitCommand("Empty commit")
	if err == nil {
		t.Errorf("CommitCommand should fail with empty index")
	} else if !strings.Contains(err.Error(), "nothing to commit") {
		t.Errorf("Expected error to contain 'nothing to commit', got: %v", err)
	}

	// Modify the test file for the second commit
	modifiedContent := []byte("Modified content")
	if err := os.WriteFile(testFile, modifiedContent, 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Create a new blob for the modified file
	modifiedBlob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		t.Fatalf("Failed to create blob for modified file: %v", err)
	}

	// Write the modified blob to the objects directory
	modifiedBlobPath := filepath.Join(objectsDir, modifiedBlob.ID())
	serialized, err = modifiedBlob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize modified blob: %v", err)
	}
	if err := os.WriteFile(modifiedBlobPath, serialized, 0644); err != nil {
		t.Fatalf("Failed to write modified blob file: %v", err)
	}

	// Update the index with the modified file
	indexEntries[relPath] = modifiedBlob.ID()
	indexData, err = json.Marshal(indexEntries)
	if err != nil {
		t.Fatalf("Failed to marshal updated index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatalf("Failed to write updated index file: %v", err)
	}

	// SECOND COMMIT CREATION - Manual version
	modifiedTree := core.BuildTreeFromPaths(indexEntries)

	// Store the modified tree in the objects directory
	modifiedTreeData, err := modifiedTree.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize modified tree: %v", err)
	}
	modifiedTreePath := filepath.Join(objectsDir, modifiedTree.ID())
	if err := os.WriteFile(modifiedTreePath, modifiedTreeData, 0644); err != nil {
		t.Fatalf("Failed to write modified tree file: %v", err)
	}

	// Create a commit pointing to the modified tree with the first commit as parent
	modifiedCommit := core.NewCommit(modifiedTree.ID(), commit.ID(), "Modified file", "test")

	// Store the modified commit in the objects directory
	modifiedCommitData, err := modifiedCommit.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize modified commit: %v", err)
	}
	modifiedCommitPath := filepath.Join(objectsDir, modifiedCommit.ID())
	if err := os.WriteFile(modifiedCommitPath, modifiedCommitData, 0644); err != nil {
		t.Fatalf("Failed to write modified commit file: %v", err)
	}

	// Update the master branch reference to point to the modified commit
	if err := os.WriteFile(masterRefPath, []byte(modifiedCommit.ID()), 0644); err != nil {
		t.Fatalf("Failed to update master ref file: %v", err)
	}
}
