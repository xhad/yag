package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhad/yag/internal/commands"
	"github.com/xhad/yag/internal/core"
	"github.com/xhad/yag/internal/storage"
	"github.com/xhad/yag/tests/testutil"
)

// TestRestoreCommand tests unstaging files with the restore command
func TestRestoreCommand(t *testing.T) {
	log := testutil.NewLogger(t)
	log.StartTest()
	defer log.EndTest()

	// Create a temporary directory for the test
	log.Section("Setting up test environment")
	startTime := time.Now()
	log.Action("Creating", "temporary directory")
	tempDir, err := os.MkdirTemp("", "yag_test_restore_*")
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
	log.Timing("Environment setup", startTime)

	// Initialize a proper repository structure
	log.Section("Creating repository structure")
	startTime = time.Now()
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
	log.Timing("Repository structure creation", startTime)

	// Create test files
	log.Section("Creating test files")
	startTime = time.Now()
	log.File("file1.txt", "Creating")
	file1 := filepath.Join(tempDir, "file1.txt")
	file1Content := []byte("File 1 content")
	if err := os.WriteFile(file1, file1Content, 0644); err != nil {
		log.Error("Failed to create file1: %v", err)
		t.Fatalf("Failed to create file1: %v", err)
	}

	log.File("file2.txt", "Creating")
	file2 := filepath.Join(tempDir, "file2.txt")
	file2Content := []byte("File 2 content")
	if err := os.WriteFile(file2, file2Content, 0644); err != nil {
		log.Error("Failed to create file2: %v", err)
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Get relative paths
	file1Rel, err := filepath.Rel(tempDir, file1)
	if err != nil {
		log.Error("Failed to get relative path for file1: %v", err)
		t.Fatalf("Failed to get relative path for file1: %v", err)
	}
	log.Info("Relative path for file1: %s", file1Rel)

	file2Rel, err := filepath.Rel(tempDir, file2)
	if err != nil {
		log.Error("Failed to get relative path for file2: %v", err)
		t.Fatalf("Failed to get relative path for file2: %v", err)
	}
	log.Info("Relative path for file2: %s", file2Rel)
	log.Timing("Test files creation", startTime)

	// Create blobs for the files
	log.Section("Creating and storing blobs")
	startTime = time.Now()
	log.Action("Creating", "blob for file1")
	blob1, err := core.NewBlobFromFile(file1)
	if err != nil {
		log.Error("Failed to create blob for file1: %v", err)
		t.Fatalf("Failed to create blob for file1: %v", err)
	}
	log.Info("Created blob for file1 with ID: %s", blob1.ID())

	log.Action("Creating", "blob for file2")
	blob2, err := core.NewBlobFromFile(file2)
	if err != nil {
		log.Error("Failed to create blob for file2: %v", err)
		t.Fatalf("Failed to create blob for file2: %v", err)
	}
	log.Info("Created blob for file2 with ID: %s", blob2.ID())

	// Write the blobs to the objects directory
	log.File(blob1.ID(), "Storing blob")
	blob1Path := filepath.Join(objectsDir, blob1.ID())
	blob1Data, err := blob1.Serialize()
	if err != nil {
		log.Error("Failed to serialize blob1: %v", err)
		t.Fatalf("Failed to serialize blob1: %v", err)
	}
	if err := os.WriteFile(blob1Path, blob1Data, 0644); err != nil {
		log.Error("Failed to write blob1 file: %v", err)
		t.Fatalf("Failed to write blob1 file: %v", err)
	}

	log.File(blob2.ID(), "Storing blob")
	blob2Path := filepath.Join(objectsDir, blob2.ID())
	blob2Data, err := blob2.Serialize()
	if err != nil {
		log.Error("Failed to serialize blob2: %v", err)
		t.Fatalf("Failed to serialize blob2: %v", err)
	}
	if err := os.WriteFile(blob2Path, blob2Data, 0644); err != nil {
		log.Error("Failed to write blob2 file: %v", err)
		t.Fatalf("Failed to write blob2 file: %v", err)
	}
	log.Timing("Blob creation and storage", startTime)

	// Create index with both files staged
	log.Section("Setting up index")
	startTime = time.Now()
	log.File(filepath.Join(yagDir, storage.IndexFile), "Creating")
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	indexEntries := map[string]string{
		file1Rel: blob1.ID(),
		file2Rel: blob2.ID(),
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
	log.Timing("Index setup", startTime)

	// Verify both files are in the index
	log.Section("Verifying initial index state")
	startTime = time.Now()
	log.Action("Reading", "index file")
	entriesBefore, err := readIndexFile(indexPath)
	if err != nil {
		log.Error("Failed to read index file: %v", err)
		t.Fatalf("Failed to read index file: %v", err)
	}

	if _, exists := entriesBefore[file1Rel]; !exists {
		log.Error("file1 should be in the index before restore")
		t.Errorf("file1 should be in the index before restore")
	} else {
		log.Success("file1 is properly staged in the index")
	}

	if _, exists := entriesBefore[file2Rel]; !exists {
		log.Error("file2 should be in the index before restore")
		t.Errorf("file2 should be in the index before restore")
	} else {
		log.Success("file2 is properly staged in the index")
	}
	log.Timing("Initial index verification", startTime)

	// Test unstaging file1 using restore --staged
	log.Section("Testing restore command")
	startTime = time.Now()
	log.Command(fmt.Sprintf("yag restore --staged %s", file1Rel))
	err = commands.RestoreCommand([]string{file1Rel}, true)
	if err != nil {
		log.Error("RestoreCommand failed: %v", err)
		t.Fatalf("RestoreCommand failed: %v", err)
	} else {
		log.Success("RestoreCommand executed successfully")
	}
	log.Timing("Restore command execution", startTime)

	// Verify file1 was unstaged and file2 remains staged
	log.Section("Verifying index state after restore")
	startTime = time.Now()
	log.Action("Reading", "index file after restore")
	entriesAfter, err := readIndexFile(indexPath)
	if err != nil {
		log.Error("Failed to read index file after restore: %v", err)
		t.Fatalf("Failed to read index file after restore: %v", err)
	}

	if _, exists := entriesAfter[file1Rel]; exists {
		log.Error("file1 should not be in the index after restore")
		t.Errorf("file1 should not be in the index after restore")
	} else {
		log.Success("file1 was correctly unstaged")
	}

	if _, exists := entriesAfter[file2Rel]; !exists {
		log.Error("file2 should still be in the index after restore")
		t.Errorf("file2 should still be in the index after restore")
	} else {
		log.Success("file2 correctly remains in the index")
	}
	log.Timing("Index verification after restore", startTime)

	// Test error cases
	log.Section("Testing error cases")
	startTime = time.Now()

	// Test error case: trying to unstage a file that's not in the index
	log.Command("yag restore --staged nonexistent.txt")
	err = commands.RestoreCommand([]string{"nonexistent.txt"}, true)
	if err == nil {
		log.Error("RestoreCommand should fail with nonexistent file")
		t.Errorf("RestoreCommand should fail with nonexistent file")
	} else {
		log.Success("Correctly failed with nonexistent file: %v", err)
	}

	// Test error case: no files specified
	log.Command("yag restore --staged")
	err = commands.RestoreCommand([]string{}, true)
	if err == nil {
		log.Error("RestoreCommand should fail with no files specified")
		t.Errorf("RestoreCommand should fail with no files specified")
	} else {
		log.Success("Correctly failed with no files specified: %v", err)
	}
	log.Timing("Error case testing", startTime)
}

// Helper function to read the index file
func readIndexFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries map[string]string
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
