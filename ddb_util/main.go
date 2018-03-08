package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joemadeus/collectd_ddb/types"
)

// Verify that their table is present and available, or create the required DDB table
//
// ATTRIBUTES:
// time     N
// host     S
// plugin   S
// interval, type, values, dstypes, dsnames
//
// KEYS:
// time-host-plugin RANGE
func main() {
	create := *flag.Bool("create", false, "create the necessary tables in DDB")
	flag.Parse()

	ddb, err := types.NewDDBPlugin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create a DDB session or instance: ", err.Error())
		os.Exit(1)
	}

	ping := ddb.Ping()
	if create {
		if ping {
			fmt.Fprintln(os.Stderr, "The table already exists. Cannot create")
			os.Exit(1)
		}

		ddb.CreateTable()
	} else if ping == false {
		fmt.Fprintln(os.Stderr, "Cannot ping the table in DDB. Does it exist?")
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, "table exists")
	os.Exit(0)
}
