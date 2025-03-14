package tests

import (
	"fmt"
	"io"
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

// TestBranchCommand tests branch creation and listing
func TestBranchCommand(t *testing.T) {
	log := testutil.NewLogger(t)
	log.StartTest()
	defer log.EndTest()

	// Use the utility function to set up a test repository with a committed file
	log.Section("Setting up test repository")
	startTime := time.Now()

	// Create a temporary directory for the test
	log.Action("Creating", "temporary directory")
	tempDir, err := os.MkdirTemp("", "yag_test_branch_*")
	if err != nil {
		log.Error("Failed to create temp directory: %v", err)
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	log.Success("Created temporary directory: %s", tempDir)

	defer func() {
		log.Action("Cleaning up", "test repository")
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

	// Initialize repository structure
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

	// Create a test file
	log.File("test_file.txt", "Creating")
	testFile := filepath.Join(tempDir, "test_file.txt")
	err = os.WriteFile(testFile, []byte("Hello, YAG!"), 0644)
	if err != nil {
		log.Error("Failed to create test file: %v", err)
		t.Fatalf("Failed to create test file: %v", err)
	}
	log.Success("Created test file with content: 'Hello, YAG!'")

	log.Timing("Environment setup", startTime)

	// Since the default GetIndexEntries implementation is just a placeholder,
	// we need to manually commit to have a valid commit to create a branch from
	log.Section("Creating initial commit")
	startTime = time.Now()

	// Create a blob from the test file
	log.Action("Creating", "blob from test file")
	log.File(testFile, "Reading")
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		log.Error("Failed to create blob: %v", err)
		t.Fatalf("Failed to create blob: %v", err)
	}
	log.Info("Created blob with ID: %s", blob.ID())

	// Store the blob
	log.Repository("Storing", "blob object")
	objectPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		log.Error("Failed to serialize blob: %v", err)
		t.Fatalf("Failed to serialize blob: %v", err)
	}
	if err := os.WriteFile(objectPath, serialized, 0644); err != nil {
		log.Error("Failed to write blob to file: %v", err)
		t.Fatalf("Failed to write blob to file: %v", err)
	}
	log.Success("Stored blob in objects directory: %s", objectPath)

	// Create the tree with the relative path to the file
	log.Action("Creating", "tree object")
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		log.Error("Failed to get relative path: %v", err)
		t.Fatalf("Failed to get relative path: %v", err)
	}
	log.Info("Relative path for test file: %s", relPath)

	treeEntries := map[string]string{
		relPath: blob.ID(),
	}
	tree := core.BuildTreeFromPaths(treeEntries)
	log.Info("Created tree with ID: %s", tree.ID())

	// Store the tree
	log.Repository("Storing", "tree object")
	treeObjectPath := filepath.Join(objectsDir, tree.ID())
	treeSerialized, err := tree.Serialize()
	if err != nil {
		log.Error("Failed to serialize tree: %v", err)
		t.Fatalf("Failed to serialize tree: %v", err)
	}
	if err := os.WriteFile(treeObjectPath, treeSerialized, 0644); err != nil {
		log.Error("Failed to write tree to file: %v", err)
		t.Fatalf("Failed to write tree to file: %v", err)
	}
	log.Success("Stored tree in objects directory: %s", treeObjectPath)

	// Create a commit
	log.Action("Creating", "commit object")
	author := "test-user"
	commitMessage := "Initial commit"
	commit := core.NewCommit(tree.ID(), "", commitMessage, author)
	log.Info("Created commit with ID: %s", commit.ID())

	// Store the commit
	log.Repository("Storing", "commit object")
	commitObjectPath := filepath.Join(objectsDir, commit.ID())
	commitSerialized, err := commit.Serialize()
	if err != nil {
		log.Error("Failed to serialize commit: %v", err)
		t.Fatalf("Failed to serialize commit: %v", err)
	}
	if err := os.WriteFile(commitObjectPath, commitSerialized, 0644); err != nil {
		log.Error("Failed to write commit to file: %v", err)
		t.Fatalf("Failed to write commit to file: %v", err)
	}
	log.Success("Stored commit in objects directory: %s", commitObjectPath)

	// Update master branch to point to the new commit
	log.Repository("Updating", "master branch reference")
	masterRefPath := filepath.Join(headsDir, storage.DefaultBranch)
	if err := os.WriteFile(masterRefPath, []byte(commit.ID()), 0644); err != nil {
		log.Error("Failed to update master ref: %v", err)
		t.Fatalf("Failed to update master ref: %v", err)
	}
	log.Success("Set master branch to commit: %s", commit.ID())

	// Update HEAD to point to master branch properly
	log.Repository("Updating", "HEAD reference")
	headContent = "ref: refs/heads/" + storage.DefaultBranch
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		log.Error("Failed to update HEAD file: %v", err)
		t.Fatalf("Failed to update HEAD file: %v", err)
	}
	log.Success("Set HEAD to reference: %s", headContent)
	log.Timing("Initial commit creation", startTime)

	// Now when we run branch commands, there's a commit in place
	log.Section("Testing branch creation")
	startTime = time.Now()

	// Test creating a new branch
	newBranch := "test-branch"
	log.Command(fmt.Sprintf("yag branch %s", newBranch))
	err = commands.BranchCommand([]string{newBranch})
	if err != nil {
		log.Error("BranchCommand failed to create branch: %v", err)
		t.Fatalf("BranchCommand failed to create branch: %v", err)
	}
	log.Success("Created branch: %s", newBranch)

	// Verify the branch file was created
	branchPath := filepath.Join(tempDir, storage.YAGDir, storage.RefsDir, storage.HeadsDir, newBranch)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		log.Error("Branch file was not created at %s", branchPath)
		t.Errorf("Branch file was not created at %s", branchPath)
	} else {
		log.Success("Branch file exists at: %s", branchPath)
	}
	log.Timing("Branch creation test", startTime)

	// Test listing branches (should include master and test-branch)
	log.Section("Testing branch listing")
	startTime = time.Now()

	// Capture stdout to verify output
	log.Info("Capturing stdout to verify branch listing output")
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	log.Command("yag branch")
	err = commands.BranchCommand([]string{})
	if err != nil {
		log.Error("BranchCommand failed to list branches: %v", err)
		t.Fatalf("BranchCommand failed to list branches: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	if err != nil {
		log.Error("Failed to read captured output: %v", err)
		t.Fatalf("Failed to read captured output: %v", err)
	}
	output := buf.String()
	log.Info("Branch listing output: %s", output)

	if !strings.Contains(output, storage.DefaultBranch) || !strings.Contains(output, newBranch) {
		log.Error("Branch listing should contain '%s' and '%s', got: %s", storage.DefaultBranch, newBranch, output)
		t.Errorf("Branch listing should contain '%s' and '%s', got: %s", storage.DefaultBranch, newBranch, output)
	} else {
		log.Success("Branch listing correctly shows both branches")
	}
	log.Timing("Branch listing test", startTime)

	// Test error case: creating a branch without any commits
	log.Section("Testing branch creation error case")
	startTime = time.Now()

	// Create a new repository without commits
	log.Action("Creating", "repository without commits")
	subDir := filepath.Join(tempDir, "no-commits")
	if err := os.Mkdir(subDir, 0755); err != nil {
		log.Error("Failed to create subdirectory: %v", err)
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	if err := os.Chdir(subDir); err != nil {
		log.Error("Failed to change to subdirectory: %v", err)
		t.Fatalf("Failed to change to subdirectory: %v", err)
	}
	log.Success("Changed to empty repository directory: %s", subDir)

	log.Command("yag init")
	err = commands.InitCommand([]string{})
	if err != nil {
		log.Error("Failed to initialize repository in subdirectory: %v", err)
		t.Fatalf("Failed to initialize repository in subdirectory: %v", err)
	}
	log.Success("Initialized empty repository")

	// Try to create a branch
	log.Command("yag branch error-branch")
	err = commands.BranchCommand([]string{"error-branch"})
	if err == nil {
		log.Error("BranchCommand should fail in a repository with no commits")
		t.Errorf("BranchCommand should fail in a repository with no commits")
	} else {
		log.Success("Correctly failed to create branch in empty repository: %v", err)
	}
	log.Timing("Branch error case test", startTime)
}
