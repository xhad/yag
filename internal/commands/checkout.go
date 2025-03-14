package commands

import (
	"fmt"
	"os"

	"github.com/xhad/yag/internal/repository"
)

// CheckoutCommand switches to the specified branch
func CheckoutCommand(branchName string) error {
	if branchName == "" {
		return fmt.Errorf("branch name is required")
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

	// Checkout the branch
	if err := repo.Checkout(branchName); err != nil {
		return err
	}

	fmt.Printf("Switched to branch '%s'\n", branchName)
	return nil
}
