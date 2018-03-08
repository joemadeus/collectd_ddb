package main

import (
	"context"
	"encoding/json"
	"time"

	"collectd.org/api"
	"collectd.org/plugin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
		TableName: aws.String("Movies"),
	}

	_, err = ddbp.ddb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

func init() {
	sess, err := session.NewSession(
		&aws.Config{Region: aws.String("us-west-2")},
	)
	if err != nil {
		panic(err)
	}

	plugin.RegisterWrite("ddb", &DDBPlugin{
		session: sess,
		ddb:     dynamodb.New(sess),
	})
}

// Ignored by the collectd daemon, but end users can verify that their tables are present
// and available by exec'ing this binary from a command line
func main() {
	// TODO
}
