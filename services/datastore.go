package services

import (
	"context"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
)

// NewDatastore -
func NewDatastore(ctx context.Context, gcpProjectID string, credJSON []byte) (*datastore.Client, error) {
	return datastore.NewClient(ctx, gcpProjectID, option.WithCredentialsJSON(credJSON))
}
