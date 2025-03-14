package commands

import (
	"fmt"
	"os"

	"github.com/xhad/yag/internal/repository"
)

// AddCommand adds files to the staging area
func AddCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nothing specified, nothing added")
	}

	// Open the repository
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	repo, err := repository.Open(path)
	if err != nil {
		return err
	}

	// Add each file
	for _, file := range args {
		if err := repo.Add(file); err != nil {
			return fmt.Errorf("failed to add '%s': %v", file, err)
		}
		fmt.Printf("Added '%s'\n", file)
	}

	return nil
}
