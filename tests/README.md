# YAG Tests

This directory contains tests for the YAG Git-like application.

## Pretty Logging

YAG tests include a pretty logging system that provides colorful, structured output during test execution. This makes it easier to understand test flow, identify failures, and track performance.

### Features

- **Colorful Output**: Different colors indicate different types of messages (success, error, etc.)
- **Structured Sections**: Tests are organized into logical sections
- **Timing Information**: Performance tracking for different parts of the test
- **Progress Indicators**: Clear visual indicators of test flow
- **Detailed Error Reports**: Better context when tests fail

### Running Tests with Pretty Logging

You can run tests with the standard Go test command:

```bash
# Run all tests with verbose output
go test -v ./tests/...

# Run a specific test file
go test -v ./tests/branch_test.go

# Run a specific test function
go test -v ./tests/branch_test.go -run TestBranchCommand
```

### Using the Pretty Test Tool

We've also included a dedicated command-line tool for running tests with pretty logging:

```bash
# Build the tool
go build -o pretty-test ./cmd/pretty-test

# Run all tests
./pretty-test

# Run a specific test
./pretty-test branch

# Show help
./pretty-test --help
```

### Using the Logger in Your Tests

You can use the pretty logger in your own tests:

```go
package tests

import (
	"testing"
	"time"
	
	"github.com/xhad/yag/tests/testutil"
)

func TestExample(t *testing.T) {
	// Create a new logger
	log := testutil.NewLogger(t)
	
	// Mark the beginning of the test
	log.StartTest()
	defer log.EndTest()
	
	// Create a section
	log.Section("Setting up test environment")
	
	// Track timing for operations
	startTime := time.Now()
	
	// Log various types of messages
	log.Action("Creating", "test directory")
	log.Success("Test directory created successfully")
	log.Info("This is some informational message")
	log.Warning("This operation might fail")
	log.Error("Something went wrong: %v", err)
	
	// Show command execution
	log.Command("yag branch new-branch")
	
	// Track timing for sections
	log.Timing("Environment setup", startTime)
}
```

See the [logger documentation](testutil/README.md) for more details on available methods. 