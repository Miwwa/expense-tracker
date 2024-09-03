# Expense Tracker

Made for https://roadmap.sh/projects/expense-tracker

A simple command-line application to manage your finances by tracking expenses.

## Installation

1. Ensure you have [Go](https://golang.org/dl/) installed on your system.
2. Clone this repository to your local machine.
3. Navigate to the project directory.
4. Run `go build` to compile the application.

```sh
git clone https://github.com/Miwwa/expense-tracker
cd expense-tracker
go build
```

## Usage

The CLI application accepts various commands with corresponding arguments.

```
Usage: expense-tracker <command> [options]

expense-tracker add --description <description> --amount <amount>
expense-tracker update --id <id> [--description <description>] [--amount <amount>]
expense-tracker delete --id <id>
expense-tracker list
expense-tracker summary [--month <number>] [--year <number>]

add --help to any command to get detailed information
```
