package commands

import (
	"fmt"
	"os"

	"github.com/xhad/yag/internal/repository"
)

// BranchCommand handles branch operations
func BranchCommand(args []string) error {
	// Open the repository
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	repo, err := repository.Open(path)
	if err != nil {
		return err
	}

	// If no branch name is provided, list all branches
	if len(args) == 0 {
		return listBranches(repo)
	}

	// Otherwise, create a new branch
	branchName := args[0]

	if err := repo.CreateBranch(branchName); err != nil {
		return err
	}

	fmt.Printf("Created branch '%s'\n", branchName)
	return nil
}

// listBranches lists all branches in the repository
func listBranches(repo *repository.Repository) error {
	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}

	currentBranch, err := getCurrentBranch(repo)
	if err != nil {
		return err
	}

	for _, branch := range branches {
		prefix := "  "
		if branch == currentBranch {
			prefix = "* "
		}
		fmt.Printf("%s%s\n", prefix, branch)
	}

	return nil
}

// getCurrentBranch gets the name of the current branch
func getCurrentBranch(repo *repository.Repository) (string, error) {
	// Use the repository's storage to get the HEAD reference
	storage := repo.GetStorage()
	if storage == nil {
		return "", fmt.Errorf("failed to get storage")
	}

	return storage.GetHead()
}
