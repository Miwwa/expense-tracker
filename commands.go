package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func HelpCmd() error {
	fmt.Print(HelpText)
	return nil
}

func AddCmd(args []string, tracker *Tracker) error {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.Usage = func() {
		fmt.Fprint(addCmd.Output(), "Usage of add:\nadd a new record to the tracker\n")
		addCmd.PrintDefaults()
	}

	description := addCmd.String("description", "", "text description, required")
	amount := addCmd.Uint("amount", 0, "money amount, required, must be more than 0")

	err := addCmd.Parse(args)
	if err != nil {
		return err
	}

	if *amount == 0 {
		addCmd.Usage()
		return errors.New("invalid amount")
	}

	if *description == "" {
		addCmd.Usage()
		return errors.New("invalid description")
	}

	record, err := tracker.Add(*description, *amount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error adding record: %v\n", err)
		return err
	}

	fmt.Printf("Expense added successfully (ID: %d)", record.Id)

	return nil
}
