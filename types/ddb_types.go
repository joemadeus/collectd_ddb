package types

import (
	"context"
	"fmt"
	"os"

	collectdApi "collectd.org/api"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	// no way to set config values from collectd, otherwise this would be there
	TableName = "collectd_ddb"
)

type DDBPlugin struct {
	session *session.Session
	ddb     *dynamodb.DynamoDB
}

// Creates a new DDBPlugin with an AWS session and DDB instance
func NewDDBPlugin() (*DDBPlugin, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "(recovered)")
		}
	}()

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &DDBPlugin{
		session: sess,
		ddb:     dynamodb.New(sess),
	}, nil
}

// Creates a table suitable for collectd plugin use
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
		TableName: aws.String(TableName),
	}

	_, err := ddbp.ddb.CreateTable(create)
	if err != nil {
		return err
	}

	return nil
}

// Returns 'true' if the table can be described -- i.e., if we're able to query it
func (ddbp *DDBPlugin) Ping() bool {
	_, err := ddbp.ddb.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(TableName)})
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not ping/describe the DDB table: ", err.Error())
	}
	return err != nil
}

// Does the write work for collectd
func (ddbp *DDBPlugin) Write(_ context.Context, vl *collectdApi.ValueList) error {
	attrValue, err := dynamodbattribute.MarshalMap(vl)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attrValue,
		TableName: aws.String(TableName),
	}

	_, err = ddbp.ddb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}
