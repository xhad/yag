package core

import (
	"fmt"
	"io/ioutil"
)

// Blob represents file content in the repository
type Blob struct {
	content []byte
	hash    string
}

// NewBlob creates a new Blob from content
func NewBlob(content []byte) *Blob {
	blob := &Blob{
		content: content,
	}

	// Calculate the hash of serialized blob
	serialized := SerializeObject(BlobType, content)
	blob.hash = CalculateHash(serialized)

	return blob
}

// NewBlobFromFile creates a new Blob from a file path
func NewBlobFromFile(path string) (*Blob, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
	}

	return NewBlob(content), nil
}

// Type returns the type of this object (implements Object interface)
func (b *Blob) Type() ObjectType {
	return BlobType
}

// ID returns the hash identifier of this blob (implements Object interface)
func (b *Blob) ID() string {
	return b.hash
}

// Content returns the content of the blob
func (b *Blob) Content() []byte {
	return b.content
}

// Size returns the size of the blob content in bytes
func (b *Blob) Size() int {
	return len(b.content)
}

// Serialize converts the blob to a byte slice for storage (implements Object interface)
func (b *Blob) Serialize() ([]byte, error) {
	return SerializeObject(BlobType, b.content), nil
}
