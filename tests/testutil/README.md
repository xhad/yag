# Test Utilities for YAG

This directory contains utility functions and structures to enhance testing for the YAG project.

## Pretty Logging

The `logger.go` file implements a pretty logging utility that helps make test output more informative and visually appealing. This makes it easier to:

- Understand test flow and execution
- Debug test failures
- Track timing information for various test operations
- Visualize the structure of test execution

### Usage Example

```go
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

### Available Logging Methods

- `StartTest()` - Marks the beginning of a test
- `EndTest()` - Marks the end of a test
- `Section(name string)` - Creates a new labeled section in the test
- `Info(msg string, args ...interface{})` - Logs an informational message
- `Success(msg string, args ...interface{})` - Logs a success message
- `Warning(msg string, args ...interface{})` - Logs a warning message
- `Error(msg string, args ...interface{})` - Logs an error message
- `Action(action, target string)` - Logs an action being performed
- `Command(cmd string)` - Logs a command being executed
- `Timing(operation string, start time.Time)` - Logs the time taken for an operation
- `Separator()` - Prints a visual separator line
- `File(path, action string)` - Logs information about a file operation
- `Repository(operation, details string)` - Logs information about a repository operation

### Color Coding

The logger uses color-coded output to differentiate between different types of messages:

- Blue: Test boundaries, sections, and timing information
- Green: Success messages
- Red: Error messages
- Yellow: Warnings and commands
- Purple: Actions
- Cyan: File operations
- White: Informational messages

### Benefits

- Makes test output more readable and structured
- Provides clear visual indications of test execution flow
- Helps identify performance bottlenecks with timing information
- Makes debugging test failures easier by providing rich context
- Standardizes logging across all tests in the project 