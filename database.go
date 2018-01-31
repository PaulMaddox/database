package database

import "errors"

var (
	// ErrInvalidCredentials is thrown if no valid AWS credentials were able to be loaded
	// from the environment variables, IAM role, or ~/.aws/credentials file
	ErrInvalidCredentials = errors.New("invalid AWS credentials")

	// ErrRecordNotFound is thrown when a Get() request is made for a non-existant record
	ErrRecordNotFound = errors.New("record not found")

	// ErrInvalidRecordID is thrown when a database record doesn't have a valid identifier (id) field
	ErrInvalidRecordID = errors.New("invalid record identifier")
)

// Table interface provides methods for interacting with a database
type Table interface {
	Get(string, interface{}) error
	Put(interface{}) error
	Update(interface{}) error
	Delete(string) error
}
