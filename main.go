package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

type DDBIdentifier struct {
	Host   string `json:"host"`
	Plugin string `json:"plugin"`
	Type   string `json:"type"`
}

// A DDBValueList is the DynamoDB's representation of a collectd ValueList
type DDBValueList struct {
	DDBIdentifier
	Time     time.Time     `json:"time"`
	Interval time.Duration `json:"interval"`
	Values   []json.Number `json:"values"`
	DSNames  []string
}

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

func (ddbp *DDBPlugin) DescribeTable() (*dynamodb.DescribeTableOutput, error) {
	return ddbp.ddb.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
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
	plugin.RegisterWrite("ddb", ddb)
}

// Ignored by the collectd daemon, but end users can verify that their tables are present
// and available by exec'ing this binary from a command line
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
	ddb, err := NewDDBPlugin()
	if err != nil {
		fmt.Errorf("could not create a DDB session or instance", err)
		os.Exit(1)
	}

}
