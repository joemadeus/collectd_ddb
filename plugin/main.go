package main

import (
	collectdPlugin "collectd.org/plugin"
	"github.com/joemadeus/collectd_ddb/types"
)

func init() {
	ddb, err := types.NewDDBPlugin()
	if err != nil {
		panic(err)
	}

	if ddb.Ping() == false {
		panic("The table needed for the DDB plugin does not exist. Run the plugin from the command line to create it")
	}

	collectdPlugin.RegisterWrite("ddb", ddb)
}

func main() {} // ignored by collectd
