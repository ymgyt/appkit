package services

import (
	"context"
	"io/ioutil"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
)

// MustDatastore -
func MustDatastore(c *datastore.Client, err error) *datastore.Client {
	if err != nil {
		panic(err)
	}
	return c
}

// NewDatastore -
func NewDatastore(ctx context.Context, gcpProjectID string, credJSON []byte) (*datastore.Client, error) {
	return datastore.NewClient(ctx, gcpProjectID, option.WithCredentialsJSON(credJSON))
}

// NewDatastoreFromFile -
func NewDatastoreFromFile(ctx context.Context, gcpProjectID string, credJSONPath string) (*datastore.Client, error) {
	b, err := ioutil.ReadFile(credJSONPath)
	if err != nil {
		return nil, err
	}
	return NewDatastore(ctx, gcpProjectID, b)
}
