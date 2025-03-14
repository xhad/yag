package commands

import (
	"fmt"
	"os"
	"sort"

	"github.com/xhad/yag/internal/repository"
)

// StatusCommand shows the working tree status
func StatusCommand(args []string) error {
	// Open the repository
	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	repo, err := repository.Open(path)
	if err != nil {
		return err
	}

	// Get repository status
	status, err := repo.Status()
	if err != nil {
		return err
	}

	// Print status header with current branch
	branch, err := repo.GetCurrentBranch()
	if err != nil {
		return err
	}
	fmt.Printf("On branch %s\n", branch)

	// Print staged files
	if len(status.Staged) > 0 {
		fmt.Println("\nChanges to be committed:")
		fmt.Println("  (use \"yag restore --staged <file>...\" to unstage)")
		fmt.Println()

		// Sort the files for consistent output
		stagedFiles := make([]string, 0, len(status.Staged))
		for file := range status.Staged {
			stagedFiles = append(stagedFiles, file)
		}
		sort.Strings(stagedFiles)

		for _, file := range stagedFiles {
			fmt.Printf("\tmodified: %s\n", file)
		}
	}

	// Print unstaged files
	if len(status.Unstaged) > 0 {
		fmt.Println("\nChanges not staged for commit:")
		fmt.Println("  (use \"yag add <file>...\" to update what will be committed)")
		fmt.Println()

		// Sort the files for consistent output
		unstagedFiles := make([]string, 0, len(status.Unstaged))
		for file := range status.Unstaged {
			unstagedFiles = append(unstagedFiles, file)
		}
		sort.Strings(unstagedFiles)

		for _, file := range unstagedFiles {
			fmt.Printf("\tmodified: %s\n", file)
		}
	}

	// Print untracked files
	if len(status.Untracked) > 0 {
		fmt.Println("\nUntracked files:")
		fmt.Println("  (use \"yag add <file>...\" to include in what will be committed)")
		fmt.Println()

		// Sort the files for consistent output
		untrackedFiles := make([]string, 0, len(status.Untracked))
		for file := range status.Untracked {
			untrackedFiles = append(untrackedFiles, file)
		}
		sort.Strings(untrackedFiles)

		for _, file := range untrackedFiles {
			fmt.Printf("\t%s\n", file)
		}
	}

	// If nothing to show, print a clean message
	if len(status.Staged) == 0 && len(status.Unstaged) == 0 && len(status.Untracked) == 0 {
		fmt.Println("\nNothing to commit, working tree clean")
	}

	return nil
}
