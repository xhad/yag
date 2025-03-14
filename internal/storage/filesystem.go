package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhad/yag/internal/core"
)

const (
	YAGDir        = ".yag"
	ObjectsDir    = "objects"
	RefsDir       = "refs"
	HeadsDir      = "heads"
	IndexFile     = "index"
	HeadFile      = "HEAD"
	DefaultBranch = "master"
)

// FileSystemStorage implements the Storage interface using the file system
type FileSystemStorage struct {
	rootPath string
}

// NewFileSystemStorage creates a new FileSystemStorage
func NewFileSystemStorage(rootPath string) *FileSystemStorage {
	return &FileSystemStorage{
		rootPath: rootPath,
	}
}

// Initialize prepares the storage for use
func (fs *FileSystemStorage) Initialize() error {
	// Create .yag directory
	if err := os.MkdirAll(filepath.Join(fs.rootPath, YAGDir), 0755); err != nil {
		return err
	}

	// Create objects directory
	if err := os.MkdirAll(filepath.Join(fs.rootPath, YAGDir, ObjectsDir), 0755); err != nil {
		return err
	}

	// Create refs/heads directory
	if err := os.MkdirAll(filepath.Join(fs.rootPath, YAGDir, RefsDir, HeadsDir), 0755); err != nil {
		return err
	}

	// Create HEAD file pointing to master branch
	headPath := filepath.Join(fs.rootPath, YAGDir, HeadFile)
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/"+DefaultBranch), 0644); err != nil {
		return err
	}

	// Create empty index file
	indexPath := filepath.Join(fs.rootPath, YAGDir, IndexFile)
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		return err
	}

	return nil
}

// objectPath returns the path to an object file
func (fs *FileSystemStorage) objectPath(hash string) string {
	return filepath.Join(fs.rootPath, YAGDir, ObjectsDir, hash)
}

// refPath returns the path to a ref file
func (fs *FileSystemStorage) refPath(name string) string {
	return filepath.Join(fs.rootPath, YAGDir, RefsDir, HeadsDir, name)
}

// StoreObject stores an object in the storage
func (fs *FileSystemStorage) StoreObject(obj core.Object) error {
	data, err := obj.Serialize()
	if err != nil {
		return err
	}

	path := fs.objectPath(obj.ID())

	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write the object to disk
	return os.WriteFile(path, data, 0644)
}

// HasObject checks if an object exists in storage
func (fs *FileSystemStorage) HasObject(hash string) (bool, error) {
	path := fs.objectPath(hash)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetObject retrieves an object from storage by its hash
func (fs *FileSystemStorage) GetObject(hash string) (core.Object, error) {
	path := fs.objectPath(hash)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	objType, objData, err := core.DeserializeObject(data)
	if err != nil {
		return nil, err
	}

	switch objType {
	case core.BlobType:
		return core.NewBlob(objData), nil
	case core.TreeType:
		// TODO: Implement Tree deserialization
		return nil, fmt.Errorf("tree deserialization not implemented")
	case core.CommitType:
		return core.DeserializeCommit(objData)
	default:
		return nil, fmt.Errorf("unknown object type: %s", objType)
	}
}

// UpdateRef updates a reference (like a branch) to point to a commit
func (fs *FileSystemStorage) UpdateRef(name string, commitHash string) error {
	refPath := fs.refPath(name)

	// Create directory if it doesn't exist
	dir := filepath.Dir(refPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(refPath, []byte(commitHash), 0644)
}

// GetRef gets the commit hash that a reference points to
func (fs *FileSystemStorage) GetRef(name string) (string, error) {
	refPath := fs.refPath(name)

	data, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("reference %s not found", name)
		}
		return "", err
	}

	return string(data), nil
}

// ListRefs lists all references (branches)
func (fs *FileSystemStorage) ListRefs() (map[string]string, error) {
	refsDir := filepath.Join(fs.rootPath, YAGDir, RefsDir, HeadsDir)

	// Read the refs directory
	files, err := os.ReadDir(refsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	refs := make(map[string]string)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(refsDir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		refs[file.Name()] = string(data)
	}

	return refs, nil
}

// GetHead returns the current HEAD reference
func (fs *FileSystemStorage) GetHead() (string, error) {
	headPath := filepath.Join(fs.rootPath, YAGDir, HeadFile)

	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	headContent := string(data)

	// If HEAD is a symbolic ref (points to a branch)
	if strings.HasPrefix(headContent, "ref: ") {
		// Extract branch name
		return strings.TrimPrefix(headContent, "ref: refs/heads/"), nil
	}

	// If HEAD is detached (points directly to a commit)
	return "", nil
}

// SetHead sets the HEAD reference
func (fs *FileSystemStorage) SetHead(ref string) error {
	headPath := filepath.Join(fs.rootPath, YAGDir, HeadFile)
	content := "ref: refs/heads/" + ref
	return os.WriteFile(headPath, []byte(content), 0644)
}

// GetHeadCommit returns the commit that HEAD points to
func (fs *FileSystemStorage) GetHeadCommit() (*core.Commit, error) {
	headPath := filepath.Join(fs.rootPath, YAGDir, HeadFile)

	data, err := os.ReadFile(headPath)
	if err != nil {
		return nil, err
	}

	headContent := string(data)

	var commitHash string

	// If HEAD is a symbolic ref (points to a branch)
	if strings.HasPrefix(headContent, "ref: ") {
		branchPath := strings.TrimPrefix(headContent, "ref: ")
		branchPath = filepath.Join(fs.rootPath, YAGDir, branchPath)

		// Read the commit hash from the branch file
		hashData, err := os.ReadFile(branchPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, nil // Branch exists but has no commits
			}
			return nil, err
		}

		commitHash = string(hashData)
	} else {
		// If HEAD is detached (points directly to a commit)
		commitHash = headContent
	}

	// Get the commit object
	obj, err := fs.GetObject(commitHash)
	if err != nil {
		return nil, err
	}

	commit, ok := obj.(*core.Commit)
	if !ok {
		return nil, fmt.Errorf("object %s is not a commit", commitHash)
	}

	return commit, nil
}

// GetIndexEntries returns the current staged files
func (fs *FileSystemStorage) GetIndexEntries() (map[string]string, error) {
	indexPath := filepath.Join(fs.rootPath, YAGDir, IndexFile)

	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	// Parse the index file JSON format
	var entries map[string]string
	if len(data) > 0 {
		if err := json.Unmarshal(data, &entries); err != nil {
			// If we can't parse the index, start with an empty map
			entries = make(map[string]string)
		}
	} else {
		entries = make(map[string]string)
	}

	return entries, nil
}

// UpdateIndex updates the staging area
func (fs *FileSystemStorage) UpdateIndex(path string, hash string) error {
	indexPath := filepath.Join(fs.rootPath, YAGDir, IndexFile)

	// Read existing index entries
	entries, err := fs.GetIndexEntries()
	if err != nil {
		return err
	}

	// Update the entry
	entries[path] = hash

	// Write back to file as JSON
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// UpdateIndexEntries updates multiple entries in the staging area at once
func (fs *FileSystemStorage) UpdateIndexEntries(entries map[string]string) error {
	indexPath := filepath.Join(fs.rootPath, YAGDir, IndexFile)

	// Write entries to file as JSON
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// ClearIndex clears the staging area
func (fs *FileSystemStorage) ClearIndex() error {
	indexPath := filepath.Join(fs.rootPath, YAGDir, IndexFile)
	return os.WriteFile(indexPath, []byte("{}"), 0644)
}
