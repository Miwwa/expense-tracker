package main

import (
	"fmt"
	"os"
)

func main() {
	err := Run(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}
}

func Run(args []string) error {
	if len(args) == 0 {
		return HelpCmd()
	}

	storage := NewStorageFromFile("./expenses.csv")
	tracker, err := NewTracker(storage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tracker: %v\n", err)
		return err
	}

	switch args[0] {
	case "help", "-h", "--help":
		return HelpCmd()
	case "add":
		return AddCmd(args[1:], tracker)
	}
	return nil
}
