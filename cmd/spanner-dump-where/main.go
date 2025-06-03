package main

import (
	"context"
	"fmt"
	"github.com/Jumpaku/spanner-dump-whare/spanner-dump"
	"log"
	"os"
	"time"
)

func main() {
	if err := Run(cli{}, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type cli struct{}

func (cli) Run(input Input) error {
	if input.ErrorMessage != "" {
		fmt.Println(GetDoc(input.Subcommand))
		panicf("Error: %s\n", input.ErrorMessage)
	}
	if input.Opt_Project == "" || input.Opt_Instance == "" || input.Opt_Database == "" {
		fmt.Println(GetDoc(input.Subcommand))
		panicf("Error: Missing parameters: -project, -instance, -database are required\n")
	}
	if len(input.Opt_From) == 0 || len(input.Opt_Where) == 0 {
		fmt.Println(GetDoc(input.Subcommand))
		panicf("Error: Missing parameters: -from and -where are required\n")
	}
	if len(input.Opt_From) != len(input.Opt_Where) {
		fmt.Println(GetDoc(input.Subcommand))
		panicf("Error: Invalid parameters: count of -from and -where must be same\n")
	}

	var timestamp *time.Time
	if input.Opt_Timestamp != "" {
		t, err := time.Parse(time.RFC3339, input.Opt_Timestamp)
		panicfIfError(err, "Error: Invalid timestamp format")
		timestamp = &t
	}

	query := make(map[string]string)
	for index, from := range input.Opt_From {
		query[from] = input.Opt_Where[index]
	}

	ctx := context.Background()
	dumper, err := spanner_dump.NewDumper(ctx,
		input.Opt_Project, input.Opt_Instance, input.Opt_Database,
		os.Stdout,
		timestamp,
		uint(input.Opt_BulkSize),
		query,
		input.Opt_Sort,
		input.Opt_Upsert,
	)
	panicfIfError(err, "Failed to create dumper")
	defer dumper.Cleanup()

	if !input.Opt_NoDdl {
		err := dumper.DumpDDLs(ctx)
		panicfIfError(err, "Failed to dump DDLs")
	}

	if !input.Opt_NoData {
		err := dumper.DumpTables(ctx)
		panicfIfError(err, "Failed to dump tables")
	}

	return nil
}

func panicfIfError(err error, format string, a ...interface{}) {
	if err != nil {
		log.Panicf(fmt.Sprintf(format, a...)+": %+v", err)
	}
}

func panicf(format string, a ...interface{}) {
	log.Panicf(fmt.Sprintf(format, a...))
}
