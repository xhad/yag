// Package core implements the fundamental data structures for YAG version control
// @title YAG Core Objects
// @author XHad
// @notice Provides the core object models and interfaces for YAG
// @dev Contains object types like blobs, trees, and commits, along with serialization utilities
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ObjectType represents the type of a YAG object
// @notice String representation of object types used for serialization and identification
type ObjectType string

const (
	// BlobType represents a file content object
	// @notice Represents a file's contents in the object database
	BlobType ObjectType = "blob"

	// TreeType represents a directory structure object
	// @notice Represents a directory structure with references to blobs and other trees
	TreeType ObjectType = "tree"

	// CommitType represents a snapshot of the repository
	// @notice Represents a point-in-time snapshot with author, message, and tree references
	CommitType ObjectType = "commit"
)

// Object represents a YAG object in the object database
// @notice Base interface implemented by all storable objects in YAG
// @dev All core objects (Blob, Tree, Commit) must implement this interface
type Object interface {
	// Type returns the type of the object
	// @return ObjectType The type of this object (blob, tree, commit)
	Type() ObjectType

	// ID returns the SHA-256 hash that identifies this object
	// @return string The unique hex-encoded SHA-256 hash of this object
	ID() string

	// Serialize converts the object to a byte slice for storage
	// @return []byte, error The serialized object data and nil on success, or nil and an error on failure
	Serialize() ([]byte, error)
}

// CalculateHash computes the SHA-256 hash of the given data
// @notice Creates a deterministic identifier for object data
// @param data The raw data to hash
// @return string The hex-encoded SHA-256 hash of the data
func CalculateHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// SerializeObject creates a byte representation of an object with its header
// @notice Prepares an object for storage by adding a type and size header
// @param objType The type of object being serialized
// @param data The raw object data
// @return []byte The complete serialized object with header
func SerializeObject(objType ObjectType, data []byte) []byte {
	header := fmt.Sprintf("%s %d\x00", objType, len(data))
	return append([]byte(header), data...)
}

// DeserializeObject extracts the object type and data from a serialized object
// @notice Parses a serialized object back into its type and data components
// @param raw The serialized object bytes
// @return ObjectType, []byte, error The object type, data, and nil on success, or empty values and an error on failure
func DeserializeObject(raw []byte) (ObjectType, []byte, error) {
	// Find the null byte that separates header from data
	nullIndex := -1
	for i, b := range raw {
		if b == 0 {
			nullIndex = i
			break
		}
	}

	if nullIndex == -1 {
		return "", nil, fmt.Errorf("invalid object format: missing null byte")
	}

	header := string(raw[:nullIndex])
	data := raw[nullIndex+1:]

	var objType ObjectType
	var size int

	_, err := fmt.Sscanf(header, "%s %d", &objType, &size)
	if err != nil {
		return "", nil, fmt.Errorf("invalid header format: %v", err)
	}

	if len(data) != size {
		return "", nil, fmt.Errorf("corrupt object: size mismatch")
	}

	return objType, data, nil
}
