package database

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

// DynamoDBAdaptor provides helper methods for DynamoDB operations
type DynamoDBAdaptor struct {
	dynamoDB *dynamodb.DynamoDB
}

// DynamoDBTable provides helper methods for DynamoDB operations
type DynamoDBTable struct {
	table    string
	dynamoDB *dynamodb.DynamoDB
}

// NewDynamoDBAdapter creates a new database adapter for interacting with DynamoDB
func NewDynamoDBAdapter() (*DynamoDBAdaptor, error) {

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	ddb := dynamodb.New(cfg)

	return &DynamoDBAdaptor{
		dynamoDB: ddb,
	}, nil

}

// Table returns a DynamoDB table that has methods to interact with rows
func (d *DynamoDBAdaptor) Table(name string) Table {
	return &DynamoDBTable{
		table:    name,
		dynamoDB: d.dynamoDB,
	}
}

// Get an item from the database
func (d *DynamoDBTable) Get(id string, v interface{}) error {

	ddbreq := d.dynamoDB.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key: map[string]dynamodb.AttributeValue{
			"id": dynamodb.AttributeValue{S: aws.String(id)},
		},
	})

	ddbresp, err := ddbreq.Send()
	if err != nil {
		return err
	}

	if len(ddbresp.Item) < 1 {
		return ErrRecordNotFound
	}

	if err := dynamodbattribute.UnmarshalMap(ddbresp.Item, v); err != nil {
		return err
	}

	return nil

}

// Put item in the database
func (d *DynamoDBTable) Put(v interface{}) error {

	// Marshall the event to DynamoDB attributes
	item, err := dynamodbattribute.MarshalMap(v)
	if err != nil {
		return err
	}

	// Write the item to DynamoDB
	ddbreq := d.dynamoDB.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      item,
	})

	if _, err := ddbreq.Send(); err != nil {
		return err
	}

	return nil

}

// Update (overwrite) an item in the database
func (d *DynamoDBTable) Update(v interface{}) error {
	return d.Put(v)
}

// Delete an item from the database
func (d *DynamoDBTable) Delete(id string) error {

	ddbreq := d.dynamoDB.DeleteItemRequest(&dynamodb.DeleteItemInput{
		TableName: aws.String(d.table),
		Key: map[string]dynamodb.AttributeValue{
			"id": dynamodb.AttributeValue{S: aws.String(id)},
		},
	})

	if _, err := ddbreq.Send(); err != nil {
		return err
	}

	return nil

}
