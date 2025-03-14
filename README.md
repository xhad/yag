# YAG - Yet Another Git

YAG is a simplified Git implementation for Duncan and John that provides basic version control capabilities.

## Features

- Initialize a repository
- Add and stage files
- Commit changes with messages
- Create branches
- Checkout branches
- Status of current branch
- Restore files from previous commits

## Design

YAG is built around a Merkle DAG structure:

- **Blob**: Represents file content, hashed for content-addressable storage
- **Tree**: Represents directories and their contents
- **Commit**: Points to a tree and contains metadata (author, message, etc.)
- **Branches**: Named pointers to specific commits

The storage is file-system based, similar to Git, using a `.yag` directory to store all objects and references.

## Getting Started

### Usage

```bash
# run the yag command
make run
# or
go run cmd/yag/main.go

# build the yag command
make build
./yag

# or
go build -o yag cmd/yag/main.go

# run the tests
make test
# or
go test ./tests
```

### Basic Usage

```bash
# Initialize a repository
./yag init

# Add files to staging
./yag add file.txt

# Commit changes
./yag commit -m "Initial commit"

# Create a new branch
./yag branch feature-branch

# Switch to a branch
./yag checkout feature-branch

# List all branches
./yag branch

# Status of current branch
./yag status

# Restore files from previous commits
./yag restore file.txt
```

## Development Decisions

1. **Language**: Go was chosen for its simplicity, strong standard library, and excellent file handling capabilities.

2. **Storage**: A simple filesystem-based storage was implemented, using content-addressable storage for objects.

3. **Merkle DAG**: The core data structure is a Merkle DAG, similar to Git, which provides cryptographic verification of content.

4. **CLI Interface**: Command-line interface was designed to be similar to Git for familiarity.

## Testing

Run the tests with:

```bash
make test
# or
go test ./tests
```

## Limitations

- No network capability (local only)
- Limited merge support
- No conflict resolution
- Simple index structure

## Future Improvements

See ROADMAP.md for future enhancements and planned features.