package main

import (
	"os"

	collectdPlugin "collectd.org/plugin"
	"github.com/joemadeus/collectd_ddb/types"
)

const (
	profileEnvVarName = "COLLECTD_DDB_PROFILE"
	regionEnvVarName  = "COLLECTD_DDB_REGION"
	tableEnvVarName   = "COLLECTD_DDB_TABLE"
)

func init() {
	profile := os.Getenv(profileEnvVarName)
	region := os.Getenv(regionEnvVarName)
	tableName := os.Getenv(tableEnvVarName)

	if profile == "" {
		collectdPlugin.Errorf("the profile env var must be set. use " + profileEnvVarName)
		panic("Please set up your env correctly. See the log for details")
	}

	if region == "" {
		collectdPlugin.Errorf("the region env var must be set. use " + regionEnvVarName)
		panic("Please set up your env correctly. See the log for details")
	}

	if tableName == "" {
		collectdPlugin.Errorf("the table env var must be set. use " + tableEnvVarName)
		panic("Please set up your env correctly. See the log for details")
	}

	ddb, err := types.NewDDBPlugin(profile, region, tableName)
	if err != nil {
		panic(err)
	}

	ping, err := ddb.Ping()
	if ping == false {
		collectdPlugin.Errorf("The table needed for the DDB plugin does not exist or can't be accessed. Run the plugin from the command line to create it. Error: %s", err.Error())
		panic("The collectd_ddb table does not exist or couldn't be accessed. See the logs for details")
	}

	collectdPlugin.RegisterWrite("ddb", ddb)
}

func main() {} // ignored by collectd
