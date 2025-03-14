package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhad/yag/internal/repository"
)

// InitCommand initializes a new repository
func InitCommand(args []string) error {
	var path string

	// If a path is provided, use it, otherwise use current directory
	if len(args) > 0 {
		path = args[0]
	} else {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
	}

	// Create the directory if it doesn't exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Initialize the repository
	_, err = repository.Init(path)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %v", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	fmt.Printf("Initialized empty YAG repository in %s\n", absPath)
	return nil
}
