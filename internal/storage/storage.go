// Package storage provides interfaces and implementations for YAG's object storage system
// @title YAG Storage System
// @author XHad
// @notice Manages low-level storage operations for the YAG version control system
// @dev The storage package is responsible for persisting all objects, references, and staging information
package storage

import (
	"github.com/xhad/yag/internal/core"
)

// Storage is an interface for storing and retrieving YAG objects
// @notice Defines the contract for any storage implementation in YAG
// @dev Any storage implementation (filesystem, database, etc.) must implement this interface
type Storage interface {
	// Initialize prepares the storage for use
	// @notice Sets up the storage system for first use
	// @return error Returns nil on success or an error if initialization fails
	Initialize() error

	// StoreObject stores an object in the storage
	// @notice Persists a core.Object to the storage backend
	// @param obj The object to store
	// @return error Returns nil on success or an error if storage fails
	StoreObject(obj core.Object) error

	// GetObject retrieves an object from storage by its hash
	// @notice Fetches and deserializes an object from storage
	// @param hash The object ID/hash to retrieve
	// @return core.Object, error Returns the object and nil on success, or nil and an error if retrieval fails
	GetObject(hash string) (core.Object, error)

	// HasObject checks if an object exists in storage
	// @notice Checks for the existence of an object without retrieving it fully
	// @param hash The object ID/hash to check
	// @return bool, error Returns true if object exists, false if not, or an error if checking fails
	HasObject(hash string) (bool, error)

	// UpdateRef updates a reference (like a branch) to point to a commit
	// @notice Changes where a named reference points to
	// @param name The name of the reference to update
	// @param commitHash The commit hash the reference should point to
	// @return error Returns nil on success or an error if the update fails
	UpdateRef(name string, commitHash string) error

	// GetRef gets the commit hash that a reference points to
	// @notice Retrieves the commit hash that a named reference points to
	// @param name The name of the reference to query
	// @return string, error Returns the commit hash and nil on success, or an empty string and an error if retrieval fails
	GetRef(name string) (string, error)

	// ListRefs lists all references (branches)
	// @notice Gets all named references and their target commit hashes
	// @return map[string]string, error Returns a map of reference names to commit hashes, or an error if listing fails
	ListRefs() (map[string]string, error)

	// GetHead returns the current HEAD reference
	// @notice Gets the current HEAD reference (usually a branch name)
	// @return string, error Returns the HEAD reference and nil on success, or an empty string and an error if retrieval fails
	GetHead() (string, error)

	// SetHead sets the HEAD reference
	// @notice Updates the HEAD reference to point to a different branch
	// @param ref The reference name to set HEAD to
	// @return error Returns nil on success or an error if the update fails
	SetHead(ref string) error

	// GetHeadCommit returns the commit that HEAD points to
	// @notice Resolves HEAD to a commit object
	// @return *core.Commit, error Returns the commit object and nil on success, or nil and an error if resolution fails
	GetHeadCommit() (*core.Commit, error)

	// GetIndexEntries returns the current staged files
	// @notice Gets all entries in the staging area (index)
	// @return map[string]string, error Returns a map of file paths to object hashes, or an error if retrieval fails
	GetIndexEntries() (map[string]string, error)

	// UpdateIndex updates the staging area
	// @notice Adds or updates a single entry in the staging area
	// @param path The file path to update in the index
	// @param hash The object hash for the file content
	// @return error Returns nil on success or an error if the update fails
	UpdateIndex(path string, hash string) error

	// UpdateIndexEntries updates multiple entries in the staging area at once
	// @notice Replaces the entire staging area with new entries
	// @param entries A map of file paths to object hashes
	// @return error Returns nil on success or an error if the update fails
	UpdateIndexEntries(entries map[string]string) error

	// ClearIndex clears the staging area
	// @notice Removes all entries from the staging area
	// @return error Returns nil on success or an error if clearing fails
	ClearIndex() error
}
