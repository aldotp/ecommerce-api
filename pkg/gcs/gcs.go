package gcs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	gcstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCS struct {
	client *gcstorage.Client
	bucket string
}

func NewCGS(cred, bucketName string) (*GCS, error) {
	client, err := gcstorage.NewClient(context.Background(), option.WithCredentialsFile(cred))
	if err != nil {
		return nil, err
	}

	return &GCS{
		client: client,
		bucket: bucketName,
	}, nil
}

var _ StorageInterface = &GCS{}

func (g *GCS) Upload(ctx context.Context, object *FileUploadObject) error {
	defer object.File.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	attr, err := g.client.Bucket(g.bucket).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("Bucket(%s).Attrs: %w", g.bucket, err)
	}

	wc := g.client.Bucket(g.bucket).Object(object.FileName).NewWriter(ctx)

	if len(attr.ACL) > 0 {
		wc.ACL = []gcstorage.ACLRule{{Entity: gcstorage.AllUsers, Role: gcstorage.RoleReader}}
	}

	if _, err := io.Copy(wc, object.File); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}

	return nil
}

func (g *GCS) Download(ctx context.Context, fileName string) (*bytes.Buffer, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	rc, err := g.client.Bucket(g.bucket).Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", fileName, err)
	}
	defer rc.Close()

	buff := new(bytes.Buffer)
	buff.ReadFrom(rc)

	return buff, nil
}

func (g *GCS) Delete(ctx context.Context, fileName string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	if err := g.client.Bucket(g.bucket).Object(fileName).Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", fileName, err)
	}

	return nil
}

func (g *GCS) GenerateUrl(fileName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucket, fileName)
}

func (g *GCS) Duplicate(ctx context.Context, oldUrl string, newFilename string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	response, err := http.Get(oldUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code: %d", response.StatusCode)
	}

	writer := g.client.Bucket(g.bucket).Object(newFilename).NewWriter(ctx)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

func (g *GCS) GeneratePresignedUrl(fileName string, expiration time.Duration) (string, error) {
	url, err := g.client.Bucket(g.bucket).SignedURL(fileName, &gcstorage.SignedURLOptions{
		Scheme:  gcstorage.SigningSchemeV4,
		Method:  "PUT",
		Expires: time.Now().Add(expiration),
	})
	if err != nil {
		return "", fmt.Errorf("SignedURL: %w", err)
	}

	return url, nil
}
