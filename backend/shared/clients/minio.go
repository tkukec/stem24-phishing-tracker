package clients

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type MinIO struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	useSSL          bool
	bucketName      string
	client          *minio.Client
}

func NewMinIO(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*MinIO, error) {

	client, err := InstantiateMinIOClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	instance := &MinIO{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		useSSL:          useSSL,
		client:          client,
		bucketName:      bucketName,
	}

	return instance, nil
}

func InstantiateMinIOClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *MinIO) Client() (*minio.Client, error) {
	if c.client == nil {
		var err error
		c.client, err = InstantiateMinIOClient(c.endpoint, c.accessKeyID, c.secretAccessKey, c.useSSL)
		if err != nil {
			return nil, err
		}
	}
	return c.client, nil
}

func (c *MinIO) createBucketIfDoesNotExist(ctx context.Context, bucketName, location string) error {
	client, err := c.Client()
	if err != nil {
		return err
	}
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		_, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists != nil {
			return errBucketExists
		}
	}
	return nil
}

type PutMinIOObjectRequest struct {
	Object            io.Reader
	ObjectName        string `json:"object-name"`
	OverwriteIfExists bool   `json:"overwrite-if-exists"`
	Size              int64  `json:"size"`
	Location          string `json:"location"`
}

func (c *MinIO) PutObject(ctx context.Context, request *PutMinIOObjectRequest) (*minio.UploadInfo, error) {
	err := c.createBucketIfDoesNotExist(ctx, c.bucketName, request.Location)
	if err != nil {
		return nil, err
	}

	if !request.OverwriteIfExists {
		if objectExists, _ := c.ObjectExists(ctx, &ObjectRequest{
			ObjectName: request.ObjectName,
		}); objectExists {
			return nil, nil
		}
	}

	client, err := c.Client()
	if err != nil {
		return nil, err
	}

	info, err := client.PutObject(ctx, c.bucketName, request.ObjectName, request.Object, request.Size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &info, nil
}

type ObjectRequest struct {
	ObjectName string `json:"object-name"`
}

func (c *MinIO) GetObject(ctx context.Context, request *ObjectRequest) (io.ReadCloser, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	object, err := client.GetObject(ctx, c.bucketName, request.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (c *MinIO) RemoveObject(ctx context.Context, request *ObjectRequest) error {
	client, err := c.Client()
	if err != nil {
		return err
	}
	err = client.RemoveObject(ctx, c.bucketName, request.ObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *MinIO) ObjectExists(ctx context.Context, request *ObjectRequest) (bool, error) {
	client, err := c.Client()
	if err != nil {
		return false, err
	}
	_, err = client.StatObject(ctx, c.bucketName, request.ObjectName, minio.StatObjectOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}
