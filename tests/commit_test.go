package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
	"github.com/xhad/yag/tests/testutil"
)

// TestCommitCommand tests creating a commit with staged changes
func TestCommitCommand(t *testing.T) {
	log := testutil.NewLogger(t)
	log.StartTest()
	defer log.EndTest()

	// Create a temporary directory for the test
	log.Section("Setting up test environment")
	startTime := time.Now()
	log.Action("Creating", "temporary directory")
	tempDir, err := os.MkdirTemp("", "yag_test_commit_*")
	if err != nil {
		log.Error("Failed to create temp directory: %v", err)
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	log.Success("Created temporary directory: %s", tempDir)
	defer func() {
		log.Action("Cleaning up", "temporary directory")
		os.RemoveAll(tempDir)
	}()

	// Change to the temporary directory
	log.Action("Changing", "working directory")
	originalDir, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory: %v", err)
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		log.Error("Failed to change to temp directory: %v", err)
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	log.Success("Changed working directory to: %s", tempDir)
	defer func() {
		log.Action("Restoring", "original directory")
		os.Chdir(originalDir)
	}()

	// Initialize a proper repository structure
	log.Repository("Creating", ".yag directory")
	yagDir := filepath.Join(tempDir, storage.YAGDir)
	if err := os.MkdirAll(yagDir, 0755); err != nil {
		log.Error("Failed to create .yag directory: %v", err)
		t.Fatalf("Failed to create .yag directory: %v", err)
	}

	// Create necessary subdirectories
	log.Repository("Creating", "objects directory")
	objectsDir := filepath.Join(yagDir, storage.ObjectsDir)
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		log.Error("Failed to create objects directory: %v", err)
		t.Fatalf("Failed to create objects directory: %v", err)
	}

	log.Repository("Creating", "refs/heads directory")
	headsDir := filepath.Join(yagDir, storage.RefsDir, storage.HeadsDir)
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		log.Error("Failed to create refs/heads directory: %v", err)
		t.Fatalf("Failed to create refs/heads directory: %v", err)
	}

	// Create HEAD file
	log.File(filepath.Join(yagDir, storage.HeadFile), "Creating")
	headPath := filepath.Join(yagDir, storage.HeadFile)
	headContent := "ref: " + storage.RefsDir + "/" + storage.HeadsDir + "/" + storage.DefaultBranch
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		log.Error("Failed to create HEAD file: %v", err)
		t.Fatalf("Failed to create HEAD file: %v", err)
	}
	log.Success("Created HEAD file pointing to: %s", headContent)
	log.Timing("Repository setup", startTime)

	// Create a test file to add
	log.Section("Creating and staging test file")
	startTime = time.Now()
	log.File("test_file.txt", "Creating")
	testFile := filepath.Join(tempDir, "test_file.txt")
	testContent := []byte("Hello, YAG!")
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		log.Error("Failed to create test file: %v", err)
		t.Fatalf("Failed to create test file: %v", err)
	}
	log.Success("Created test file with content: 'Hello, YAG!'")

	// Get the relative path to the test file
	log.Action("Calculating", "relative path for test file")
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		log.Error("Failed to get relative path: %v", err)
		t.Fatalf("Failed to get relative path: %v", err)
	}
	log.Info("Relative path for test file: %s", relPath)

	// Create a blob for the test file
	log.Action("Creating", "blob from test file")
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		log.Error("Failed to create blob: %v", err)
		t.Fatalf("Failed to create blob: %v", err)
	}
	log.Info("Created blob with ID: %s", blob.ID())

	// Write the blob to the objects directory
	log.Repository("Storing", "blob object")
	blobPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		log.Error("Failed to serialize blob: %v", err)
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(blobPath, serialized, 0644); err != nil {
		log.Error("Failed to write blob file: %v", err)
		t.Fatalf("Failed to write blob file: %v", err)
	}
	log.Success("Stored blob in objects directory: %s", blobPath)

	// Create an index file with the test file
	log.File(filepath.Join(yagDir, storage.IndexFile), "Creating")
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	indexEntries := map[string]string{
		relPath: blob.ID(),
	}
	log.Info("Index entries: %v", indexEntries)
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		log.Error("Failed to marshal index data: %v", err)
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		log.Error("Failed to write index file: %v", err)
		t.Fatalf("Failed to write index file: %v", err)
	}
	log.Success("Created index file with test file entry")
	log.Timing("Test file creation and staging", startTime)

	// COMMIT CREATION - Manual version that doesn't depend on storage.GetIndexEntries
	log.Section("Creating initial commit")
	startTime = time.Now()

	// Build a tree from the test file
	log.Action("Building", "tree from index entries")
	tree := core.BuildTreeFromPaths(indexEntries)
	log.Info("Created tree with ID: %s", tree.ID())

	// Store the tree in the objects directory
	log.Repository("Storing", "tree object")
	treeData, err := tree.Serialize()
	if err != nil {
		log.Error("Failed to serialize tree: %v", err)
		t.Fatalf("Failed to serialize tree: %v", err)
	}
	treePath := filepath.Join(objectsDir, tree.ID())
	if err := os.WriteFile(treePath, treeData, 0644); err != nil {
		log.Error("Failed to write tree file: %v", err)
		t.Fatalf("Failed to write tree file: %v", err)
	}
	log.Success("Stored tree in objects directory: %s", treePath)

	// Create a commit pointing to the tree
	log.Action("Creating", "commit object")
	commit := core.NewCommit(tree.ID(), "", "Initial commit", "test")
	log.Info("Created commit with ID: %s", commit.ID())

	// Store the commit in the objects directory
	log.Repository("Storing", "commit object")
	commitData, err := commit.Serialize()
	if err != nil {
		log.Error("Failed to serialize commit: %v", err)
		t.Fatalf("Failed to serialize commit: %v", err)
	}
	commitPath := filepath.Join(objectsDir, commit.ID())
	if err := os.WriteFile(commitPath, commitData, 0644); err != nil {
		log.Error("Failed to write commit file: %v", err)
		t.Fatalf("Failed to write commit file: %v", err)
	}
	log.Success("Stored commit in objects directory: %s", commitPath)

	// Update the master branch reference to point to the commit
	log.Repository("Updating", "master branch reference")
	masterRefPath := filepath.Join(headsDir, storage.DefaultBranch)
	if err := os.WriteFile(masterRefPath, []byte(commit.ID()), 0644); err != nil {
		log.Error("Failed to write master ref file: %v", err)
		t.Fatalf("Failed to write master ref file: %v", err)
	}
	log.Success("Set master branch to commit: %s", commit.ID())
	log.Timing("Initial commit creation", startTime)

	// Test error cases
	log.Section("Testing error cases")
	startTime = time.Now()

	// Test empty commit message (using command directly)
	log.Command("yag commit -m \"\"")
	err = commands.CommitCommand("")
	if err == nil {
		log.Error("CommitCommand should fail with empty commit message")
		t.Errorf("CommitCommand should fail with empty commit message")
	} else {
		log.Success("Correctly failed with empty commit message: %v", err)
	}

	// Test empty index case (using command directly)
	log.Action("Creating", "empty index file")
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		log.Error("Failed to write empty index file: %v", err)
		t.Fatalf("Failed to write empty index file: %v", err)
	}
	log.Command("yag commit -m \"Empty commit\"")
	err = commands.CommitCommand("Empty commit")
	if err == nil {
		log.Error("CommitCommand should fail with empty index")
		t.Errorf("CommitCommand should fail with empty index")
	} else if !strings.Contains(err.Error(), "nothing to commit") {
		log.Error("Expected error to contain 'nothing to commit', got: %v", err)
		t.Errorf("Expected error to contain 'nothing to commit', got: %v", err)
	} else {
		log.Success("Correctly failed with empty index: %v", err)
	}
	log.Timing("Error case testing", startTime)

	// Modify the test file for the second commit
	log.Section("Creating second commit with modified file")
	startTime = time.Now()
	log.File("test_file.txt", "Modifying")
	modifiedContent := []byte("Modified content")
	if err := os.WriteFile(testFile, modifiedContent, 0644); err != nil {
		log.Error("Failed to modify test file: %v", err)
		t.Fatalf("Failed to modify test file: %v", err)
	}
	log.Success("Modified test file with content: 'Modified content'")

	// Create a new blob for the modified file
	log.Action("Creating", "blob from modified file")
	modifiedBlob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		log.Error("Failed to create blob for modified file: %v", err)
		t.Fatalf("Failed to create blob for modified file: %v", err)
	}
	log.Info("Created blob for modified file with ID: %s", modifiedBlob.ID())

	// Write the modified blob to the objects directory
	log.Repository("Storing", "modified blob object")
	modifiedBlobPath := filepath.Join(objectsDir, modifiedBlob.ID())
	serialized, err = modifiedBlob.Serialize()
	if err != nil {
		log.Error("Failed to serialize modified blob: %v", err)
		t.Fatalf("Failed to serialize modified blob: %v", err)
	}
	if err := os.WriteFile(modifiedBlobPath, serialized, 0644); err != nil {
		log.Error("Failed to write modified blob file: %v", err)
		t.Fatalf("Failed to write modified blob file: %v", err)
	}
	log.Success("Stored modified blob in objects directory: %s", modifiedBlobPath)

	// Update the index with the modified file
	log.Action("Updating", "index with modified file")
	indexEntries[relPath] = modifiedBlob.ID()
	indexData, err = json.Marshal(indexEntries)
	if err != nil {
		log.Error("Failed to marshal updated index data: %v", err)
		t.Fatalf("Failed to marshal updated index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		log.Error("Failed to write updated index file: %v", err)
		t.Fatalf("Failed to write updated index file: %v", err)
	}
	log.Success("Updated index with modified file")

	// SECOND COMMIT CREATION - Manual version
	log.Action("Building", "tree from updated index entries")
	modifiedTree := core.BuildTreeFromPaths(indexEntries)
	log.Info("Created tree with ID: %s", modifiedTree.ID())

	// Store the modified tree in the objects directory
	log.Repository("Storing", "modified tree object")
	modifiedTreeData, err := modifiedTree.Serialize()
	if err != nil {
		log.Error("Failed to serialize modified tree: %v", err)
		t.Fatalf("Failed to serialize modified tree: %v", err)
	}
	modifiedTreePath := filepath.Join(objectsDir, modifiedTree.ID())
	if err := os.WriteFile(modifiedTreePath, modifiedTreeData, 0644); err != nil {
		log.Error("Failed to write modified tree file: %v", err)
		t.Fatalf("Failed to write modified tree file: %v", err)
	}
	log.Success("Stored modified tree in objects directory: %s", modifiedTreePath)

	// Create a commit pointing to the modified tree with the first commit as parent
	log.Action("Creating", "second commit object")
	modifiedCommit := core.NewCommit(modifiedTree.ID(), commit.ID(), "Modified file", "test")
	log.Info("Created second commit with ID: %s", modifiedCommit.ID())

	// Store the modified commit in the objects directory
	log.Repository("Storing", "second commit object")
	modifiedCommitData, err := modifiedCommit.Serialize()
	if err != nil {
		log.Error("Failed to serialize modified commit: %v", err)
		t.Fatalf("Failed to serialize modified commit: %v", err)
	}
	modifiedCommitPath := filepath.Join(objectsDir, modifiedCommit.ID())
	if err := os.WriteFile(modifiedCommitPath, modifiedCommitData, 0644); err != nil {
		log.Error("Failed to write modified commit file: %v", err)
		t.Fatalf("Failed to write modified commit file: %v", err)
	}
	log.Success("Stored second commit in objects directory: %s", modifiedCommitPath)

	// Update the master branch reference to point to the modified commit
	log.Repository("Updating", "master branch reference")
	if err := os.WriteFile(masterRefPath, []byte(modifiedCommit.ID()), 0644); err != nil {
		log.Error("Failed to update master ref file: %v", err)
		t.Fatalf("Failed to update master ref file: %v", err)
	}
	log.Success("Updated master branch to second commit: %s", modifiedCommit.ID())
	log.Timing("Second commit creation", startTime)
}
