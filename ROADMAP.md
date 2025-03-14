# YAG Roadmap

This document outlines the planned future enhancements for the YAG project.

## Near-term Improvements

### Storage Enhancements
- [ ] Add compression for stored objects
- [ ] Implement object packing for better storage efficiency
- [ ] Improve index format to track file metadata
- [ ] Add garbage collection for unreferenced objects

### Core Functionality
- [ ] Implement diff functionality between commits
- [ ] Add basic merge capabilities (fast-forward)
- [ ] Support for tagging specific commits
- [ ] Implement stashing of working directory changes

### User Experience
- [ ] Add status command to show working tree status
- [ ] Implement log command to view commit history
- [ ] Add help command with detailed documentation
- [ ] Improve error messages with suggestions

## Mid-term Goals

### Advanced Features
- [ ] Implement basic conflict resolution
- [ ] Support for interactive rebasing
- [ ] Add cherry-pick functionality
- [ ] Implement three-way merge algorithm
- [ ] Support for signing commits

### Performance
- [ ] Optimize object storage for large repositories
- [ ] Improve performance for repositories with many files
- [ ] Add caching for frequently accessed objects
- [ ] Benchmark and optimize critical paths

### User Interface
- [ ] Add basic TUI (Text User Interface)
- [ ] Implement pager for long outputs
- [ ] Add color support for terminal output
- [ ] Interactive staging (partial file commits)

## Long-term Vision

### Advanced Capabilities
- [ ] Remote repository support (push/pull)
- [ ] Basic networking protocol
- [ ] Authentication and authorization
- [ ] Support for hooks
- [ ] Plugin architecture for extensions

### Ecosystem
- [ ] Hosting service integration
- [ ] CI/CD integration
- [ ] Web interface
- [ ] Language-specific tools and integrations

## Expansion Ideas

- [ ] Distributed storage backends
- [ ] Support for large binary files
- [ ] Advanced visualization tools
- [ ] Custom merge drivers
- [ ] Git compatibility layer 