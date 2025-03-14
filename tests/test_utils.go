package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
	"github.com/xhad/yag/tests/testutil"
)

// SetupTestRepo creates a basic repository with a file and commits it
// Returns the temporary directory path, file path, and any error
func SetupTestRepo() (string, string, error) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_setup_*")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	// Save current directory
	originalDir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		return "", "", fmt.Errorf("failed to change to temp directory: %v", err)
	}

	// Initialize repo
	err = initTestRepo(tempDir)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to initialize test repo: %v", err)
	}

	// Create, add, and commit a file
	testFile := filepath.Join(tempDir, "test_file.txt")
	err = os.WriteFile(testFile, []byte("Hello, YAG!"), 0644)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to create test file: %v", err)
	}

	// Manually update the index since the implementation is incomplete
	indexPath := filepath.Join(tempDir, storage.YAGDir, "index")

	// Create blob and store it
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to create blob: %v", err)
	}

	// Store blob in the objects directory
	objectsDir := filepath.Join(tempDir, storage.YAGDir, "objects")
	objectPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to serialize blob: %v", err)
	}
	err = os.WriteFile(objectPath, serialized, 0644)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to write blob to file: %v", err)
	}

	// Create a test index file with the proper format
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to get relative path: %v", err)
	}

	// Create index entries
	indexEntries := map[string]string{
		relPath: blob.ID(),
	}

	// Write index to file as JSON
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to marshal index data: %v", err)
	}

	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to write index file: %v", err)
	}

	// Return to original directory
	os.Chdir(originalDir)

	return tempDir, testFile, nil
}

// SetupTestRepoWithLogger creates a basic repository with a file and commits it
// Uses the pretty logger for detailed test output
// Returns the temporary directory path, file path, and any error
func SetupTestRepoWithLogger(t *testing.T, log *testutil.Logger) (string, string, error) {
	startTime := time.Now()
	log.Action("Creating", "temporary directory")

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_setup_*")
	if err != nil {
		log.Error("Failed to create temp directory: %v", err)
		return "", "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	log.Success("Created temporary directory: %s", tempDir)

	// Save current directory
	log.Action("Saving", "current directory")
	originalDir, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current directory: %v", err)
		return "", "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Change to the temporary directory
	log.Action("Changing", "to test directory")
	if err := os.Chdir(tempDir); err != nil {
		log.Error("Failed to change to temp directory: %v", err)
		return "", "", fmt.Errorf("failed to change to temp directory: %v", err)
	}
	log.Success("Changed to test directory: %s", tempDir)

	// Initialize repo
	log.Action("Initializing", "test repository")
	err = initTestRepoWithLogger(tempDir, log)
	if err != nil {
		log.Error("Failed to initialize test repo: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to initialize test repo: %v", err)
	}
	log.Success("Initialized test repository")

	// Create a test file
	log.File("test_file.txt", "Creating")
	testFile := filepath.Join(tempDir, "test_file.txt")
	err = os.WriteFile(testFile, []byte("Hello, YAG!"), 0644)
	if err != nil {
		log.Error("Failed to create test file: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to create test file: %v", err)
	}
	log.Success("Created test file with content: 'Hello, YAG!'")

	// Manually update the index
	indexPath := filepath.Join(tempDir, storage.YAGDir, storage.IndexFile)
	log.File(indexPath, "Preparing to update")

	// Create blob and store it
	log.Action("Creating", "blob from test file")
	blob, err := core.NewBlobFromFile(testFile)
	if err != nil {
		log.Error("Failed to create blob: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to create blob: %v", err)
	}
	log.Info("Created blob with ID: %s", blob.ID())

	// Store blob in the objects directory
	log.Repository("Storing", "blob object")
	objectsDir := filepath.Join(tempDir, storage.YAGDir, storage.ObjectsDir)
	objectPath := filepath.Join(objectsDir, blob.ID())
	serialized, err := blob.Serialize()
	if err != nil {
		log.Error("Failed to serialize blob: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to serialize blob: %v", err)
	}

	err = os.WriteFile(objectPath, serialized, 0644)
	if err != nil {
		log.Error("Failed to write blob to file: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to write blob to file: %v", err)
	}
	log.Success("Stored blob in objects directory: %s", objectPath)

	// Create a test index file with the proper format
	log.Action("Calculating", "relative path for test file")
	relPath, err := filepath.Rel(tempDir, testFile)
	if err != nil {
		log.Error("Failed to get relative path: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to get relative path: %v", err)
	}
	log.Info("Relative path: %s", relPath)

	// Create index entries
	log.Action("Creating", "index entries")
	indexEntries := map[string]string{
		relPath: blob.ID(),
	}
	log.Info("Index entries: %v", indexEntries)

	// Write index to file as JSON
	log.Action("Serializing", "index entries to JSON")
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		log.Error("Failed to marshal index data: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to marshal index data: %v", err)
	}

	log.File(indexPath, "Writing")
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		log.Error("Failed to write index file: %v", err)
		os.Chdir(originalDir)
		return "", "", fmt.Errorf("failed to write index file: %v", err)
	}
	log.Success("Updated index file with test file entry")

	// Return to original directory
	log.Action("Restoring", "original directory")
	os.Chdir(originalDir)
	log.Timing("Test repository setup", startTime)

	return tempDir, testFile, nil
}

// initTestRepo initializes a repository for testing
func initTestRepo(path string) error {
	// Create .yag directory
	yagDir := filepath.Join(path, storage.YAGDir)
	if err := os.MkdirAll(yagDir, 0755); err != nil {
		return fmt.Errorf("failed to create .yag directory: %v", err)
	}

	// Create objects directory
	objectsDir := filepath.Join(yagDir, storage.ObjectsDir)
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	// Create refs/heads directory
	headsDir := filepath.Join(yagDir, storage.RefsDir, storage.HeadsDir)
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		return fmt.Errorf("failed to create refs/heads directory: %v", err)
	}

	// Create HEAD file that points to master branch with the correct format
	headPath := filepath.Join(yagDir, storage.HeadFile)
	headContent := fmt.Sprintf("ref: %s/%s/%s", storage.RefsDir, storage.HeadsDir, storage.DefaultBranch)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	// Create an empty index file
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		return fmt.Errorf("failed to create index file: %v", err)
	}

	return nil
}

// initTestRepoWithLogger initializes a repository for testing with logging
func initTestRepoWithLogger(path string, log *testutil.Logger) error {
	// Create .yag directory
	log.Repository("Creating", ".yag directory")
	yagDir := filepath.Join(path, storage.YAGDir)
	if err := os.MkdirAll(yagDir, 0755); err != nil {
		log.Error("Failed to create .yag directory: %v", err)
		return fmt.Errorf("failed to create .yag directory: %v", err)
	}

	// Create objects directory
	log.Repository("Creating", "objects directory")
	objectsDir := filepath.Join(yagDir, storage.ObjectsDir)
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		log.Error("Failed to create objects directory: %v", err)
		return fmt.Errorf("failed to create objects directory: %v", err)
	}

	// Create refs/heads directory
	log.Repository("Creating", "refs/heads directory")
	headsDir := filepath.Join(yagDir, storage.RefsDir, storage.HeadsDir)
	if err := os.MkdirAll(headsDir, 0755); err != nil {
		log.Error("Failed to create refs/heads directory: %v", err)
		return fmt.Errorf("failed to create refs/heads directory: %v", err)
	}

	// Create HEAD file
	log.File(filepath.Join(yagDir, storage.HeadFile), "Creating")
	headPath := filepath.Join(yagDir, storage.HeadFile)
	headContent := fmt.Sprintf("ref: %s/%s/%s", storage.RefsDir, storage.HeadsDir, storage.DefaultBranch)
	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		log.Error("Failed to create HEAD file: %v", err)
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}
	log.Success("Created HEAD file pointing to: %s", headContent)

	// Create an empty index file
	log.File(filepath.Join(yagDir, storage.IndexFile), "Creating")
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		log.Error("Failed to create index file: %v", err)
		return fmt.Errorf("failed to create index file: %v", err)
	}
	log.Success("Created empty index file")

	return nil
}
