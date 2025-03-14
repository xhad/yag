package commands

import (
	"fmt"
	"os"

	"github.com/xhad/yag/internal/repository"
)

// RestoreCommand handles restoring files from the staging area
// @notice Removes files from the staging area when used with the --staged flag
// @dev Currently only supports unstaging files; restoring working tree changes is not implemented
// @param args The file paths to be unstaged
// @param staged Boolean flag indicating whether to unstage files (true) or restore working tree (false)
// @return error Returns nil on success or an error if the operation fails
func RestoreCommand(args []string, staged bool) error {
	if len(args) == 0 {
		return fmt.Errorf("nothing specified, nothing restored")
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

	// Check if we're unstaging files
	if staged {
		for _, file := range args {
			if err := repo.Unstage(file); err != nil {
				return fmt.Errorf("failed to unstage '%s': %v", file, err)
			}
			fmt.Printf("Unstaged changes for '%s'\n", file)
		}
		return nil
	}

	// TODO: Implement restoring working tree changes
	// (Discarding local modifications)
	return fmt.Errorf("restoring working tree changes is not yet implemented")
}
