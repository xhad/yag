package tests

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/xhad/yag/internal/storage"
)

// MockFileSystemStorage enhances the original FileSystemStorage for testing
type MockFileSystemStorage struct {
	storage.Storage
	rootPath string
}

// CreateMockStorage creates a mock storage instance that reads the actual index file
func CreateMockStorage(path string) storage.Storage {
	// Create a new mock storage
	return &MockFileSystemStorage{
		rootPath: path,
	}
}

// GetIndexEntries is a testable version that reads the actual index file
func (fs *MockFileSystemStorage) GetIndexEntries() (map[string]string, error) {
	indexPath := filepath.Join(fs.rootPath, storage.YAGDir, "index")

	// Check if the index file exists
	_, err := os.Stat(indexPath)
	if os.IsNotExist(err) {
		return make(map[string]string), nil
	}
	if err != nil {
		return nil, err
	}

	// Read the index file
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON
	if len(data) > 0 {
		var entries map[string]string
		if err := json.Unmarshal(data, &entries); err != nil {
			return nil, err
		}
		return entries, nil
	}

	return make(map[string]string), nil
}

// UpdateTestIndex simplifies updating the index file in tests
func UpdateTestIndex(repoPath, filePath, hash string) error {
	// Get relative path to the repository root
	relPath, err := filepath.Rel(repoPath, filePath)
	if err != nil {
		return err
	}

	indexPath := filepath.Join(repoPath, storage.YAGDir, "index")

	// Read existing index
	var entries map[string]string
	data, err := os.ReadFile(indexPath)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &entries); err == nil {
			// Successfully parsed existing data
		} else {
			// Start fresh if we can't parse
			entries = make(map[string]string)
		}
	} else {
		// Start fresh if we can't read
		entries = make(map[string]string)
	}

	// Update the entry
	entries[relPath] = hash

	// Write back to the index file
	indexData, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, indexData, 0644)
}
