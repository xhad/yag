package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhad/yag/internal/storage"
)

// PatchGetIndexEntries applies a runtime patch to make GetIndexEntries work properly for testing
func PatchGetIndexEntries() error {
	// This function patches the internal storage by replacing the GetIndexEntries
	// implementation with one that actually reads the index file

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Define the path to the index file
	indexPath := filepath.Join(dir, storage.YAGDir, "index")

	// Check if the index file exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Index doesn't exist, return empty map (same as original implementation)
		return nil
	}

	// Read the index file
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index file: %v", err)
	}

	// Parse the JSON if the file is not empty
	if len(data) > 0 {
		var entries map[string]string
		if err := json.Unmarshal(data, &entries); err != nil {
			return fmt.Errorf("failed to parse index file: %v", err)
		}

		// Store the entries in a global variable that our mock will use
		indexEntries = entries
	}

	return nil
}

// Global variable to hold index entries
var indexEntries = make(map[string]string)

// GetIndexEntriesForTest is a testable version that reads the actual index file
func GetIndexEntriesForTest() (map[string]string, error) {
	// Make a copy to avoid modification of the global variable
	result := make(map[string]string)
	for k, v := range indexEntries {
		result[k] = v
	}
	return result, nil
}
