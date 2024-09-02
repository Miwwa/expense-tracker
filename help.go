package main

const HelpText = `Usage: expense-tracker <command> [options]

expense-tracker add --description <description> --amount <amount>
expense-tracker update --id <id> [--description <description>] [--amount <amount>]
expense-tracker delete --id <id>
expense-tracker list
expense-tracker summary [--month <number>] [--year <number>]

add --help to any command to get detailed information
`

const AddHelpText = `expense-tracker add --description <description> --amount <amount>

add a new record to the tracker
--description <string>    text description, required
--amount <number>         money amount, required, must be more than 0
`

const UpdateHelpText = `expense-tracker update --id <id> [--description <description>] [--amount <amount>]

set new description and/or amount to record with specified id, at least one optional parameter must be specified
--id <id>                 record id, required
--description <string>    text description, optional
--amount <number>         money amount, optional, must be more than 0
`

const DeleteHelpText = `expense-tracker delete --id <id>

delete record with specified id
`

const ListHelpText = `expense-tracker list

list all tracked records
`

const SummaryHelpText = `expense-tracker summary [--month [number]] [--year [number]]

show total expenses for all time, can set optional parameters to show total expenses for specified period
--month             show total expenses for the current month
--month <number>    show total expenses for the specified month (from 1 to 12)
--year              show total expenses for the current year
--year  <number>    show total expenses for the specified year
`
