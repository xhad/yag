package testutil

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunAllTests runs all the tests in the project with verbose output
// to showcase the pretty logging.
func RunAllTests() {
	fmt.Println("Running all tests with pretty logging...")
	fmt.Println()

	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	// Find project root (look for go.mod)
	projectRoot := findProjectRoot(wd)
	if projectRoot == "" {
		fmt.Println("Could not find project root (no go.mod found)")
		return
	}

	// Run tests with verbose output
	cmd := exec.Command("go", "test", "-v", "./tests/...")
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
	}
}

// findProjectRoot looks for a go.mod file to determine the project root
func findProjectRoot(startDir string) string {
	// Check if go.mod exists in current directory
	_, err := os.Stat(filepath.Join(startDir, "go.mod"))
	if err == nil {
		return startDir
	}

	// Try parent directory if possible
	parent := filepath.Dir(startDir)
	if parent == startDir {
		// Reached root without finding go.mod
		return ""
	}

	return findProjectRoot(parent)
}

// Run runs a specific test with pretty logging
func Run(testName string) {
	if testName == "" {
		fmt.Println("No test name provided")
		return
	}

	fmt.Printf("Running test %s with pretty logging...\n\n", testName)

	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		return
	}

	// Find project root (look for go.mod)
	projectRoot := findProjectRoot(wd)
	if projectRoot == "" {
		fmt.Println("Could not find project root (no go.mod found)")
		return
	}

	// Find the test file
	testFiles, err := findTestFiles(projectRoot, testName)
	if err != nil {
		fmt.Printf("Error finding test files: %v\n", err)
		return
	}

	if len(testFiles) == 0 {
		fmt.Printf("No test files found matching %s\n", testName)
		return
	}

	// Run the test(s)
	args := []string{"test", "-v"}
	args = append(args, testFiles...)
	cmd := exec.Command("go", args...)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
	}
}

// findTestFiles finds test files matching the given test name
func findTestFiles(projectRoot, testName string) ([]string, error) {
	var testFiles []string

	// Remove Test prefix if it's there to make it more flexible
	searchName := strings.TrimPrefix(testName, "Test")

	err := filepath.Walk(filepath.Join(projectRoot, "tests"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			// If the file name contains the search name (case insensitive)
			if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(searchName)) {
				relPath, err := filepath.Rel(projectRoot, path)
				if err != nil {
					return err
				}
				testFiles = append(testFiles, "./"+relPath)
			}
		}

		return nil
	})

	return testFiles, err
}

// Helper function to parse command-line arguments and run tests
func init() {
	// If the package was imported just for the logger, don't do anything
	if len(os.Args) <= 1 || !strings.HasSuffix(os.Args[0], "testutil") {
		return
	}

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 || args[0] == "all" {
		RunAllTests()
	} else {
		Run(args[0])
	}

	os.Exit(0)
}
