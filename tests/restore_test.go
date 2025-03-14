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

// TestRestoreCommand tests unstaging files with the restore command
func TestRestoreCommand(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "yag_test_restore_*")
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

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file1Content := []byte("File 1 content")
	if err := os.WriteFile(file1, file1Content, 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	file2 := filepath.Join(tempDir, "file2.txt")
	file2Content := []byte("File 2 content")
	if err := os.WriteFile(file2, file2Content, 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Get relative paths
	file1Rel, err := filepath.Rel(tempDir, file1)
	if err != nil {
		t.Fatalf("Failed to get relative path for file1: %v", err)
	}

	file2Rel, err := filepath.Rel(tempDir, file2)
	if err != nil {
		t.Fatalf("Failed to get relative path for file2: %v", err)
	}

	// Create blobs for the files
	blob1, err := core.NewBlobFromFile(file1)
	if err != nil {
		t.Fatalf("Failed to create blob for file1: %v", err)
	}

	blob2, err := core.NewBlobFromFile(file2)
	if err != nil {
		t.Fatalf("Failed to create blob for file2: %v", err)
	}

	// Write the blobs to the objects directory
	blob1Path := filepath.Join(objectsDir, blob1.ID())
	blob1Data, err := blob1.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob1: %v", err)
	}
	if err := os.WriteFile(blob1Path, blob1Data, 0644); err != nil {
		t.Fatalf("Failed to write blob1 file: %v", err)
	}

	blob2Path := filepath.Join(objectsDir, blob2.ID())
	blob2Data, err := blob2.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob2: %v", err)
	}
	if err := os.WriteFile(blob2Path, blob2Data, 0644); err != nil {
		t.Fatalf("Failed to write blob2 file: %v", err)
	}

	// Create index with both files staged
	indexPath := filepath.Join(yagDir, storage.IndexFile)
	indexEntries := map[string]string{
		file1Rel: blob1.ID(),
		file2Rel: blob2.ID(),
	}
	indexData, err := json.Marshal(indexEntries)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// Verify both files are in the index
	entriesBefore, err := readIndexFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file: %v", err)
	}

	if _, exists := entriesBefore[file1Rel]; !exists {
		t.Errorf("file1 should be in the index before restore")
	}

	if _, exists := entriesBefore[file2Rel]; !exists {
		t.Errorf("file2 should be in the index before restore")
	}

	// Test unstaging file1 using restore --staged
	err = commands.RestoreCommand([]string{file1Rel}, true)
	if err != nil {
		t.Fatalf("RestoreCommand failed: %v", err)
	}

	// Verify file1 was unstaged and file2 remains staged
	entriesAfter, err := readIndexFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file after restore: %v", err)
	}

	if _, exists := entriesAfter[file1Rel]; exists {
		t.Errorf("file1 should not be in the index after restore")
	}

	if _, exists := entriesAfter[file2Rel]; !exists {
		t.Errorf("file2 should still be in the index after restore")
	}

	// Test error cases
	// Test error case: trying to unstage a file that's not in the index
	err = commands.RestoreCommand([]string{"nonexistent.txt"}, true)
	if err == nil {
		t.Errorf("RestoreCommand should fail with nonexistent file")
	}

	// Test error case: no files specified
	err = commands.RestoreCommand([]string{}, true)
	if err == nil {
		t.Errorf("RestoreCommand should fail with no files specified")
	}
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
