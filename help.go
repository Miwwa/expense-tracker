package main

const HelpText = `Usage: expense-tracker <command> [options]

expense-tracker add --description <description> --amount <amount>
expense-tracker update --id <id> [--description <description>] [--amount <amount>]
expense-tracker delete --id <id>
expense-tracker list
expense-tracker summary [--month <number>] [--year <number>]

add --help to any command to get detailed information
`
