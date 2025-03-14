package commands

import (
	"fmt"
	"os"

	"github.com/xhad/yag/internal/repository"
)

// CommitCommand creates a new commit with the current staged changes
func CommitCommand(message string) error {
	if message == "" {
		return fmt.Errorf("aborting commit due to empty commit message")
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

	// Create the commit
	commitID, err := repo.Commit(message)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] %s\n", commitID[:8], message)
	return nil
}
