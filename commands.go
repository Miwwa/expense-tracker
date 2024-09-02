package main

import "fmt"

type CmdError struct {
	message string
	usage   string
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("%s\nUsage: task-cli %s", e.message, e.usage)
}

func InvalidUsageError(usage string) error {
	return &CmdError{message: "invalid arguments", usage: usage}
}

func HelpCmd() (string, error) {
	return HelpText, nil
}
