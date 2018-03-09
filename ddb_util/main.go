package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joemadeus/collectd_ddb/types"
)

// Verify that their table is present and available, or create the required DDB table
func main() {
	profile := flag.String("profile", "collectd_ddb", "the AWS profile name to use when creating or describing tables")
	region := flag.String("region", "us-east-1", "the AWS region in which to store DDB data")
	tableName := flag.String("tablename", "collectd_ddb", "the name of the DDB table in which to store data")
	create := flag.Bool("create", false, "create the necessary tables in DDB")
	flag.Parse()

	fmt.Fprintf(os.Stdout, "profile %s in region %s with table %s\n", *profile, *region, *tableName)
	ddb, err := types.NewDDBPlugin(*region, *profile, *tableName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create a DDB session or instance: ", err.Error())
		os.Exit(1)
	}

	ping, err := ddb.Ping()
	if *create {
		if ping {
			fmt.Fprintln(os.Stderr, "The table already exists. Cannot create")
			os.Exit(1)
		}

		err = ddb.CreateTable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the table: %v\n", err.Error())
		} else {
			fmt.Fprintf(os.Stdout, "Created table %s with read and write provisioned throughput of 10\n")
		}

	} else if ping == false {
		fmt.Fprintf(os.Stderr, "Cannot ping the table in DDB: %v\n", err.Error())
		os.Exit(1)

	} else {
		fmt.Fprint(os.Stdout, "table exists\n")
	}

	os.Exit(0)
}
