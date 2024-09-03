package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
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

func ListCmd(_ []string, tracker *Tracker) error {
	fmt.Println("ID\tDate\t\tDescription\t\tAmount")
	for _, record := range tracker.GetAll() {
		fmt.Printf("%d\t%s\t%s\t%d\n", record.Id, record.CreatedAt.Format(time.DateOnly), record.Description, record.Amount)
	}
	return nil
}

func SummaryCmd(args []string, tracker *Tracker) error {
	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryCmd.Usage = func() {
		fmt.Fprint(summaryCmd.Output(), "Usage of summary:\nshow total expenses for all time, can set optional parameters to show total expenses for specified period\n")
		summaryCmd.PrintDefaults()
	}

	month := summaryCmd.Int("month", int(time.Now().Month()), "show total expenses for the specified month (1-12)")
	year := summaryCmd.Int("year", time.Now().Year(), "show total expenses for the specified month (1-12)")

	err := summaryCmd.Parse(args)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		fmt.Printf("Total expenses: %d", tracker.GetSummary())
		return nil
	}

	isMonthPassed := isFlagPassed(summaryCmd, "month")
	isYearPassed := isFlagPassed(summaryCmd, "year")

	if isMonthPassed && (*month < 1 || *month > 12) {
		summaryCmd.Usage()
		return errors.New("invalid month")
	}
	if isYearPassed && (*year < 1970 || *year > 9999) {
		summaryCmd.Usage()
		return errors.New("invalid year")
	}

	var sum uint
	if isYearPassed && !isMonthPassed {
		sum = tracker.GetSummaryByYear(*year)
	} else {
		sum = tracker.GetSummaryByMonth(time.Month(*month), *year)
	}
	fmt.Printf("Total expenses: %d", sum)
	return nil
}

func isFlagPassed(flags *flag.FlagSet, name string) bool {
	found := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
