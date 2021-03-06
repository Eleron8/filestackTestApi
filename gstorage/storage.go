package gstorage

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type StorageData struct {
	ProjectID     string
	BucketName    string
	ctx           context.Context
	StorageClient *storage.Client
}

func NewStorageData(ctx context.Context, projectId, bucketname, jsonPath string) (*StorageData, error) {
	storageclient, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		return nil, err
	}
	return &StorageData{
		ProjectID:     projectId,
		BucketName:    bucketname,
		ctx:           ctx,
		StorageClient: storageclient,
	}, nil
}

func (s *StorageData) GetWriter() io.Writer {
	bucket := s.StorageClient.Bucket(s.BucketName)
	wc := bucket.Object("archive.zip").NewWriter(s.ctx)

	return wc
}
