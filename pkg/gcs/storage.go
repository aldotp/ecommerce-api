package gcs

import (
	"bytes"
	"context"
)

type StorageInterface interface {
	Upload(ctx context.Context, file *FileUploadObject) error
	Download(ctx context.Context, fileName string) (*bytes.Buffer, error)
	Delete(ctx context.Context, fileName string) error
}
