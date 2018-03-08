package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"collectd.org/api"
	"collectd.org/plugin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	// no way to set config values from collectd, otherwise this would be there
	tableName = "collectd_ddb"
)

type DDBPlugin struct {
	session *session.Session
	ddb     *dynamodb.DynamoDB
}

func (ddbp *DDBPlugin) Write(_ context.Context, vl *api.ValueList) error {
	attrValue, err := dynamodbattribute.MarshalMap(vl)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attrValue,
		TableName: aws.String(tableName),
	}

	_, err = ddbp.ddb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

func (ddbp *DDBPlugin) CreateTable() error {
	create := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("time"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("host"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("plugin"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("time-host-plugin"),
				KeyType:       aws.String("RANGE"),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := ddbp.ddb.CreateTable(create)
	if err != nil {
		return err
	}

	return nil
}

func (ddbp *DDBPlugin) Ping() bool {
	_, err := ddbp.ddb.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
	if err != nil {
		fmt.Fprint(os.Stderr, "could not ping/describe the DDB table", err)
	}
	return err != nil
}

func NewDDBPlugin() (*DDBPlugin, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &DDBPlugin{
		session: sess,
		ddb:     dynamodb.New(sess),
	}, nil
}

func init() {
	ddb, err := NewDDBPlugin()
	if err != nil {
		panic(err)
	}

	tableExists := ddb.Ping()
	if tableExists == false {
		panic("The table needed for the DDB plugin does not exist. Run the plugin from the command line to create it")
	}

	plugin.RegisterWrite("ddb", ddb)
}

// Ignored by the collectd daemon, but end users can verify that their table is present
// and available -- or create the table -- by exec'ing this binary from a command line
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
	var create = flag.Bool("create", false, "create the necessary tables in DDB")
	flag.Parse()

	ddb, err := NewDDBPlugin()
	if err != nil {
		fmt.Fprint(os.Stderr, "could not create a DDB session or instance", err)
		os.Exit(1)
	}

	ping := ddb.Ping()
	if *create {
		if ping {
			fmt.Fprint(os.Stderr, "The table already exists. Cannot create")
			os.Exit(1)
		}

		ddb.CreateTable()
	}

	fmt.Fprint(os.Stdout, "table exists")
	os.Exit(0)
}
