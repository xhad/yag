package tests

import (
	"bytes"
	"testing"

	"github.com/xhad/yag/internal/core"
)

// TestBlobCreationAndSerialization tests creating a blob and serializing it
func TestBlobCreationAndSerialization(t *testing.T) {
	// Create test content
	content := []byte("Hello, YAG!")

	// Create a blob
	blob := core.NewBlob(content)

	// Verify the blob has the correct type
	if blob.Type() != core.BlobType {
		t.Errorf("Expected blob type to be %s, got %s", core.BlobType, blob.Type())
	}

	// Verify the content is preserved
	if !bytes.Equal(blob.Content(), content) {
		t.Errorf("Blob content doesn't match original content")
	}

	// Test serialization
	serialized, err := blob.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize blob: %v", err)
	}

	// Extract the data portion from serialized blob
	objType, objData, err := core.DeserializeObject(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize blob: %v", err)
	}

	// Verify the object type
	if objType != core.BlobType {
		t.Errorf("Expected deserialized type to be %s, got %s", core.BlobType, objType)
	}

	// Verify the content
	if !bytes.Equal(objData, content) {
		t.Errorf("Deserialized content doesn't match original content")
	}

	// Ensure hash generation works
	if blob.ID() == "" {
		t.Errorf("Blob ID (hash) is empty")
	}
}
