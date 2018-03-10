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

type DDBValueList struct {
	*collectdApi.ValueList
	HostPlugin string `json:"host-plugin"`
}

type DDBPlugin struct {
	TableName string
	session   *session.Session
	ddb       *dynamodb.DynamoDB
}

// Creates a new DDBPlugin with an AWS session and DDB instance
func NewDDBPlugin(awsRegion string, awsProfileName string, ddbTableName string) (*DDBPlugin, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "(recovered)")
		}
	}()

	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			Profile: awsProfileName,
			Config: aws.Config{
				Region:                        aws.String(awsRegion),
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
		}))

	return &DDBPlugin{
		TableName: ddbTableName,
		session:   sess,
		ddb:       dynamodb.New(sess),
	}, nil
}

// Creates a table suitable for collectd plugin use
//
// ATTRIBUTES: host-plugin (S) & time (N)
// KEYS: host-plugin (HASH) & time (RANGE)
func (ddbp *DDBPlugin) CreateTable() error {
	create := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("host-plugin"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("time"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("host-plugin"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("time"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(ddbp.TableName),
	}

	_, err := ddbp.ddb.CreateTable(create)
	if err != nil {
		return err
	}

	return nil
}

// Returns 'true' if the table can be described -- i.e., if we're able to query it
func (ddbp *DDBPlugin) Ping() (bool, error) {
	_, err := ddbp.ddb.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(ddbp.TableName)})
	return err == nil, err
}

// Does the write work for collectd
func (ddbp *DDBPlugin) Write(_ context.Context, vl *collectdApi.ValueList) error {
	hpBytes := make([]byte, 0)
	posit := 0

	posit += copy(hpBytes[posit:], vl.Host)
	posit += copy(hpBytes[posit:], "-")
	posit += copy(hpBytes[posit:], vl.Plugin)

	ddbValueList := DDBValueList{
		vl,
		string(hpBytes),
	}

	attrValue, err := dynamodbattribute.MarshalMap(ddbValueList)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      attrValue,
		TableName: aws.String(ddbp.TableName),
	}

	_, err = ddbp.ddb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}
