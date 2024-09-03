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

func UpdateCmd(args []string, tracker *Tracker) error {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateCmd.Usage = func() {
		fmt.Fprint(updateCmd.Output(), "Usage of update:\nset new description and/or amount to record with specified id, at least one optional parameter must be specified\n")
		updateCmd.PrintDefaults()
	}

	id := updateCmd.Uint("id", InvalidId, "record ID, required")
	description := updateCmd.String("description", "", "new text description")
	amount := updateCmd.Uint("amount", DoNotUpdateAmount, "new money amount")

	err := updateCmd.Parse(args)
	if err != nil {
		return err
	}

	if *id == InvalidId {
		updateCmd.Usage()
		return errors.New("invalid ID")
	}

	if *description == "" && *amount == DoNotUpdateAmount {
		updateCmd.Usage()
		return errors.New("required description or amount")
	}

	record, err := tracker.Update(RecordId(*id), *description, *amount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error updating record: %v\n", err)
		return err
	}

	fmt.Printf("Record updated successfully (ID: %d)", record.Id)

	return nil
}

func DeleteCmd(args []string, tracker *Tracker) error {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteCmd.Usage = func() {
		fmt.Fprint(deleteCmd.Output(), "Usage of delete:\ndelete record with specified id\n")
		deleteCmd.PrintDefaults()
	}

	id := deleteCmd.Uint("id", InvalidId, "record ID, required")

	err := deleteCmd.Parse(args)
	if err != nil {
		return err
	}

	if *id == InvalidId {
		deleteCmd.Usage()
		return errors.New("invalid ID")
	}

	err = tracker.Delete(RecordId(*id))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error deleting record: %v\n", err)
		return err
	}

	fmt.Printf("Record deleted successfully (ID: %d)", *id)

	return nil
}
