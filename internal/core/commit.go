package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// CommitData contains the data for a commit
type CommitData struct {
	TreeHash   string    // Hash of the tree this commit points to
	ParentHash string    // Hash of the parent commit (empty for root commit)
	Message    string    // Commit message
	Author     string    // Author of the commit
	Timestamp  time.Time // When the commit was created
}

// Commit represents a commit in the repository
type Commit struct {
	data CommitData
	hash string
}

// NewCommit creates a new Commit
func NewCommit(treeHash, parentHash, message, author string) *Commit {
	commit := &Commit{
		data: CommitData{
			TreeHash:   treeHash,
			ParentHash: parentHash,
			Message:    message,
			Author:     author,
			Timestamp:  time.Now(),
		},
	}

	// Calculate hash
	data, _ := commit.Serialize()
	commit.hash = CalculateHash(data)

	return commit
}

// Type returns the type of this object (implements Object interface)
func (c *Commit) Type() ObjectType {
	return CommitType
}

// ID returns the hash identifier of this commit (implements Object interface)
func (c *Commit) ID() string {
	return c.hash
}

// TreeHash returns the hash of the tree this commit points to
func (c *Commit) TreeHash() string {
	return c.data.TreeHash
}

// ParentHash returns the hash of the parent commit
func (c *Commit) ParentHash() string {
	return c.data.ParentHash
}

// Message returns the commit message
func (c *Commit) Message() string {
	return c.data.Message
}

// Author returns the commit author
func (c *Commit) Author() string {
	return c.data.Author
}

// Timestamp returns when the commit was created
func (c *Commit) Timestamp() time.Time {
	return c.data.Timestamp
}

// Serialize converts the commit to a byte slice for storage (implements Object interface)
func (c *Commit) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(c.data); err != nil {
		return nil, fmt.Errorf("failed to encode commit: %v", err)
	}

	return SerializeObject(CommitType, buf.Bytes()), nil
}

// DeserializeCommit creates a Commit from serialized data
func DeserializeCommit(data []byte) (*Commit, error) {
	var commitData CommitData

	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&commitData); err != nil {
		return nil, fmt.Errorf("failed to decode commit: %v", err)
	}

	commit := &Commit{
		data: commitData,
	}

	// Calculate hash
	serialized, _ := commit.Serialize()
	commit.hash = CalculateHash(serialized)

	return commit, nil
}
