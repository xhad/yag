package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xhad/yag/internal/commands"
)

func main() {
	// Define command line subcommands
	if len(os.Args) < 2 {
		fmt.Println("Usage: yag <command> [<args>]")
		fmt.Println("Available commands: init, add, commit, branch, checkout, status, restore")
		os.Exit(1)
	}

	// Parse the command
	command := os.Args[1]

	// Remove the command from the arguments
	os.Args = append(os.Args[:1], os.Args[2:]...)

	var err error

	// Handle commands
	switch command {
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		initCmd.Parse(os.Args[1:])
		err = commands.InitCommand(initCmd.Args())

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		addCmd.Parse(os.Args[1:])
		if addCmd.NArg() == 0 {
			fmt.Println("Usage: yag add <file1> [<file2> ...]")
			os.Exit(1)
		}
		err = commands.AddCommand(addCmd.Args())

	case "commit":
		commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
		message := commitCmd.String("m", "", "Commit message")
		commitCmd.Parse(os.Args[1:])
		if *message == "" {
			fmt.Println("Error: Commit message is required (-m flag)")
			os.Exit(1)
		}
		err = commands.CommitCommand(*message)

	case "branch":
		branchCmd := flag.NewFlagSet("branch", flag.ExitOnError)
		branchCmd.Parse(os.Args[1:])
		err = commands.BranchCommand(branchCmd.Args())

	case "checkout":
		checkoutCmd := flag.NewFlagSet("checkout", flag.ExitOnError)
		checkoutCmd.Parse(os.Args[1:])
		if checkoutCmd.NArg() == 0 {
			fmt.Println("Usage: yag checkout <branch>")
			os.Exit(1)
		}
		err = commands.CheckoutCommand(checkoutCmd.Arg(0))

	case "status":
		statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
		statusCmd.Parse(os.Args[1:])
		err = commands.StatusCommand(statusCmd.Args())

	case "restore":
		restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
		staged := restoreCmd.Bool("staged", false, "Restore staged changes (unstage files)")
		restoreCmd.Parse(os.Args[1:])

		if restoreCmd.NArg() == 0 {
			fmt.Println("Usage: yag restore [--staged] <file1> [<file2> ...]")
			os.Exit(1)
		}

		err = commands.RestoreCommand(restoreCmd.Args(), *staged)

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: init, add, commit, branch, checkout, status, restore")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
