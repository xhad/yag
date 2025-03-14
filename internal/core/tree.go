package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"sort"
)

// TreeEntry represents an entry in a tree (file or directory)
type TreeEntry struct {
	Name string    // Name of the file or directory
	Hash string    // Hash of the object
	Mode EntryMode // Mode (file or directory)
}

// EntryMode represents the mode of a tree entry
type EntryMode int

const (
	ModeFile EntryMode = 0100644
	ModeDir  EntryMode = 0040000
)

// Tree represents a directory in the repository
type Tree struct {
	entries []*TreeEntry
	hash    string
}

// NewTree creates a new Tree with no entries
func NewTree() *Tree {
	return &Tree{
		entries: []*TreeEntry{},
	}
}

// Type returns the type of this object (implements Object interface)
func (t *Tree) Type() ObjectType {
	return TreeType
}

// ID returns the hash identifier of this tree (implements Object interface)
func (t *Tree) ID() string {
	if t.hash == "" {
		// Calculate hash on demand if not already done
		data, _ := t.Serialize()
		t.hash = CalculateHash(data)
	}
	return t.hash
}

// AddEntry adds a new entry to the tree
func (t *Tree) AddEntry(name string, hash string, mode EntryMode) {
	t.entries = append(t.entries, &TreeEntry{
		Name: name,
		Hash: hash,
		Mode: mode,
	})

	// Reset hash since tree has changed
	t.hash = ""
}

// AddFile adds a file entry to the tree
func (t *Tree) AddFile(name string, hash string) {
	t.AddEntry(name, hash, ModeFile)
}

// AddDirectory adds a directory entry to the tree
func (t *Tree) AddDirectory(name string, hash string) {
	t.AddEntry(name, hash, ModeDir)
}

// GetEntries returns all entries in the tree
func (t *Tree) GetEntries() []*TreeEntry {
	// Return a sorted copy of entries (sorted by name)
	sorted := make([]*TreeEntry, len(t.entries))
	copy(sorted, t.entries)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}

// Serialize converts the tree to a byte slice for storage (implements Object interface)
func (t *Tree) Serialize() ([]byte, error) {
	// Sort entries by name for consistent hashing
	entries := t.GetEntries()

	// Encode entries using gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entries); err != nil {
		return nil, fmt.Errorf("failed to encode tree entries: %v", err)
	}

	return SerializeObject(TreeType, buf.Bytes()), nil
}

// BuildTreeFromPaths constructs a tree structure from a set of paths and their blob hashes
func BuildTreeFromPaths(paths map[string]string) *Tree {
	// Group files by directory
	dirMap := make(map[string]map[string]string)

	for path, hash := range paths {
		dir, file := filepath.Split(path)
		dir = filepath.Clean(dir)

		if dir == "." {
			dir = ""
		}

		if _, exists := dirMap[dir]; !exists {
			dirMap[dir] = make(map[string]string)
		}

		dirMap[dir][file] = hash
	}

	// Build trees from the bottom up
	treeMap := make(map[string]string)

	// Process directory by directory
	var processDirs func(string) string
	processDirs = func(dir string) string {
		// Check if this directory was already processed
		if hash, exists := treeMap[dir]; exists {
			return hash
		}

		tree := NewTree()

		// Add all files in this directory
		for file, hash := range dirMap[dir] {
			tree.AddFile(file, hash)
		}

		// Add all subdirectories
		for otherDir := range dirMap {
			if otherDir != dir && filepath.Dir(otherDir) == dir {
				subTreeHash := processDirs(otherDir)
				tree.AddDirectory(filepath.Base(otherDir), subTreeHash)
			}
		}

		// Store tree hash for reuse
		treeMap[dir] = tree.ID()

		return tree.ID()
	}

	// Start with root directory
	rootTree := NewTree()
	for dir := range dirMap {
		if filepath.Dir(dir) == "." {
			subTreeHash := processDirs(dir)
			rootTree.AddDirectory(dir, subTreeHash)
		}
	}

	// Add root-level files
	for file, hash := range dirMap[""] {
		rootTree.AddFile(file, hash)
	}

	return rootTree
}
