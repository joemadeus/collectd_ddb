# collectd_ddb
A (somewhat crudely assembled) collectd plugin that writes to AWS's DynamoDB. This was written to support data acquisition in a system that sits in the hot sun, the snow, the wind and the rain to collect weather and other information from the environment.

### Building
You must have the `collectd` sources available in order to build this plugin (I compile it with v5.8.) Change the path in the build script to point to the `collectd` sources and run it, creating two artifacts:
* `ddb`: a command line tool that, when run without options, will tell you if your credentials are good and if the machine you're running it on can access the table needed by the plugin. Running the same tool with `--create` will create the needed table. There are other flags there you can use to specify the AWS profile, region and the name of the DDB table.
* `ddb.so`: this is the plugin for `collectd`. Add it to your `collectd` plugin configuration. This plugin needs three environment variables set:
** `COLLECTD_DDB_PROFILE`: the name of the AWS profile to use
** `COLLECTD_DDB_REGION`: the region in which the DynamoDB table was created
** `COLLECTD_DDB_TABLE`: the name of the DynamoDB table into which `collectd` data is written. This one is optional and defaults to `collectd_ddb`

The table's primary key is composite, host/plugin (hash key) plus time (range key). Note that the user running the command line tool or `collectd` must have a `.aws` directory in their home dir, and that must contain `config` and `credentials` files. See the AWS docs for details on these files.

(Yes, I know that setting env variables is a lousy way to go about this but the golang plugin doesn't have any way to handle configuration at the moment.)

Note that this is still a work in progress. (I don't think it even works yet, tbh.)

<a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-sa/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/">Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License</a>.
