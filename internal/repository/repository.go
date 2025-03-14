// Package repository implements high-level repository operations for YAG
// @title YAG Repository Management
// @author XHad
// @notice Provides functionality to manage YAG repositories, including operations like add, commit, branch
// @dev Uses the storage package for persistence and core package for object representations
package repository

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/xhad/yag/internal/core"

	"github.com/xhad/yag/internal/storage"
)

// RepositoryStatus represents the status of files in the repository
// @notice Contains the categorized status of files in the repository for status command
// @dev Maps file paths to booleans for efficient lookups
type RepositoryStatus struct {
	Staged    map[string]bool // Files staged for commit
	Unstaged  map[string]bool // Files modified but not staged
	Untracked map[string]bool // Files not tracked by YAG
}

// Repository represents a YAG repository
// @notice The main structure for interacting with a YAG repository
// @dev Encapsulates storage implementation and provides high-level operations
type Repository struct {
	storage storage.Storage
	path    string
}

// Init initializes a new repository at the given path
// @notice Creates a new YAG repository in the specified directory
// @param path The directory path where the repository should be created
// @return *Repository, error The initialized repository and nil on success, or nil and an error on failure
func Init(path string) (*Repository, error) {
	// Create a new repository
	repo := &Repository{
		path: path,
	}

	// Create filesystem storage
	repo.storage = storage.NewFileSystemStorage(path)

	// Initialize the storage
	if err := repo.storage.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %v", err)
	}

	return repo, nil
}

// Open opens an existing repository at the given path
func Open(path string) (*Repository, error) {
	// Check if .yag directory exists
	yagDir := filepath.Join(path, storage.YAGDir)
	_, err := os.Stat(yagDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("not a yag repository (or any of the parent directories): %s", path)
	}
	if err != nil {
		return nil, err
	}

	// Create a new repository
	repo := &Repository{
		path: path,
	}

	// Create filesystem storage
	repo.storage = storage.NewFileSystemStorage(path)

	return repo, nil
}

// Add adds a file to the staging area
func (r *Repository) Add(filePath string) error {
	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// Check if file exists
	fi, err := os.Stat(absPath)
	if err != nil {
		return err
	}

	// If path is a directory, add all files in the directory
	if fi.IsDir() {
		return r.addDirectory(absPath)
	}

	// Add a single file
	return r.addFile(absPath)
}

// addFile adds a single file to the staging area
func (r *Repository) addFile(absPath string) error {
	// Create blob from file
	blob, err := core.NewBlobFromFile(absPath)
	if err != nil {
		return err
	}

	// Store blob in object database
	if err := r.storage.StoreObject(blob); err != nil {
		return err
	}

	// Get relative path to repository root
	relPath, err := filepath.Rel(r.path, absPath)
	if err != nil {
		return err
	}

	// Add to index
	return r.storage.UpdateIndex(relPath, blob.ID())
}

// addDirectory recursively adds all files in a directory
func (r *Repository) addDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .yag directory
		if filepath.Base(path) == storage.YAGDir {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Add file
		return r.addFile(path)
	})
}

// Commit creates a new commit with the current staged files
func (r *Repository) Commit(message string) (string, error) {
	// Get current staged files
	stagedFiles, err := r.storage.GetIndexEntries()
	if err != nil {
		return "", err
	}

	if len(stagedFiles) == 0 {
		return "", fmt.Errorf("nothing to commit, working tree clean")
	}

	// Build a tree from staged files
	tree := core.BuildTreeFromPaths(stagedFiles)

	// Store tree in object database
	if err := r.storage.StoreObject(tree); err != nil {
		return "", err
	}

	// Get parent commit hash
	var parentHash string
	headCommit, err := r.storage.GetHeadCommit()
	if err == nil && headCommit != nil {
		// Get commit ID via core.Object interface to satisfy the linter
		var obj core.Object = headCommit
		parentHash = obj.ID()
	}

	// Get author information
	currentUser, err := user.Current()
	var author string
	if err == nil {
		author = currentUser.Username
	} else {
		author = "unknown"
	}

	// Create commit
	commit := core.NewCommit(tree.ID(), parentHash, message, author)

	// Store commit in object database
	if err := r.storage.StoreObject(commit); err != nil {
		return "", err
	}

	// Update current branch to point to the new commit
	head, err := r.storage.GetHead()
	if err != nil {
		return "", err
	}

	if err := r.storage.UpdateRef(head, commit.ID()); err != nil {
		return "", err
	}

	// Clear index
	if err := r.storage.ClearIndex(); err != nil {
		return "", err
	}

	return commit.ID(), nil
}

// CreateBranch creates a new branch pointing to the current HEAD
func (r *Repository) CreateBranch(name string) error {
	// Get current HEAD commit
	headCommit, err := r.storage.GetHeadCommit()
	if err != nil {
		return err
	}

	if headCommit == nil {
		return fmt.Errorf("cannot create branch '%s': you must create at least one commit first", name)
	}

	// Update the branch reference
	return r.storage.UpdateRef(name, headCommit.ID())
}

// ListBranches lists all branches in the repository
func (r *Repository) ListBranches() ([]string, error) {
	refs, err := r.storage.ListRefs()
	if err != nil {
		return nil, err
	}

	branches := make([]string, 0, len(refs))
	for branch := range refs {
		branches = append(branches, branch)
	}

	return branches, nil
}

// Checkout switches to the specified branch
func (r *Repository) Checkout(branchName string) error {
	// Check if branch exists
	_, err := r.storage.GetRef(branchName)
	if err != nil {
		return fmt.Errorf("branch '%s' does not exist", branchName)
	}

	// Update HEAD to point to the branch
	return r.storage.SetHead(branchName)
}

// GetStorage returns the repository's storage
func (r *Repository) GetStorage() storage.Storage {
	return r.storage
}

// GetCurrentBranch returns the name of the current branch
func (r *Repository) GetCurrentBranch() (string, error) {
	return r.storage.GetHead()
}

// Status returns the status of files in the repository
func (r *Repository) Status() (*RepositoryStatus, error) {
	// Initialize status
	status := &RepositoryStatus{
		Staged:    make(map[string]bool),
		Unstaged:  make(map[string]bool),
		Untracked: make(map[string]bool),
	}

	// Get staged files from index
	indexEntries, err := r.storage.GetIndexEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get index entries: %v", err)
	}

	// Get all files in the workspace
	workspaceFiles := make(map[string]bool)
	if err := filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .yag directory
		if info.IsDir() && filepath.Base(path) == storage.YAGDir {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(r.path, path)
		if err != nil {
			return err
		}

		workspaceFiles[relPath] = true
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk workspace: %v", err)
	}

	// Compare workspace files with index
	for file := range workspaceFiles {
		_, inIndex := indexEntries[file]
		if inIndex {
			// File is in index, check if it's been modified
			filePath := filepath.Join(r.path, file)
			blob, err := core.NewBlobFromFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create blob from file: %v", err)
			}

			// If the hash is different, file is unstaged
			if blob.ID() != indexEntries[file] {
				status.Unstaged[file] = true
			}
		} else {
			// File is not in index, it's untracked
			status.Untracked[file] = true
		}
	}

	// Add all staged files
	for file := range indexEntries {
		status.Staged[file] = true
	}

	return status, nil
}

// Unstage removes a file from the staging area
// @notice Removes a file's changes from the staging area (index)
// @dev Gets current index entries, converts the path to a relative path, removes the entry, and updates the index
// @param filePath The path to the file to unstage (can be absolute or relative)
// @return error Returns nil on success or an error if unstaging fails
func (r *Repository) Unstage(filePath string) error {
	// Get the index entries
	indexEntries, err := r.storage.GetIndexEntries()
	if err != nil {
		return fmt.Errorf("failed to get index entries: %v", err)
	}

	// Get absolute path and convert to relative path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	relPath, err := filepath.Rel(r.path, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %v", err)
	}

	// Check if file is in the index
	if _, exists := indexEntries[relPath]; !exists {
		return fmt.Errorf("pathspec '%s' did not match any file in the index", filePath)
	}

	// Remove the entry from the index
	delete(indexEntries, relPath)

	// Update the index file
	return r.storage.UpdateIndexEntries(indexEntries)
}
