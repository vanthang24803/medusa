package upload

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ProviderType defines the storage provider (MinIO, Cloudflare R2, AWS S3)
type ProviderType string

const (
	ProviderMinIO ProviderType = "minio"
	ProviderR2    ProviderType = "r2"
	ProviderS3    ProviderType = "s3"
)

// Uploader defines a common interface for all file upload services.
type Uploader interface {
	Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
}

// Config contains all configuration parameters for the Uploader.
type Config struct {
	Provider        ProviderType
	Endpoint        string // MinIO: localhost:9000 | R2: <id>.r2.cloudflarestorage.com | S3: s3.amazonaws.com
	Region          string // E.g.: "us-east-1", "auto" (for R2)
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	PublicURL       string // Public CDN URL (optional)
}

// S3CompatibleUploader implements Uploader interface for all S3-compatible storage.
type S3CompatibleUploader struct {
	client    *minio.Client
	endpoint  string
	useSSL    bool
	publicURL string
}

// NewUploader is a Factory function to create the corresponding uploader without changing business code.
func NewUploader(cfg Config) (Uploader, error) {
	// Handle default values for each Provider
	switch cfg.Provider {
	case ProviderR2:
		cfg.UseSSL = true
		if cfg.Region == "" {
			cfg.Region = "auto"
		}
	case ProviderS3:
		cfg.UseSSL = true
		if cfg.Endpoint == "" && cfg.Region != "" {
			cfg.Endpoint = "s3." + cfg.Region + ".amazonaws.com"
		}
	case ProviderMinIO:
		// MinIO defaults to the configuration directly provided by the user
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}

	return &S3CompatibleUploader{
		client:    client,
		endpoint:  cfg.Endpoint,
		useSSL:    cfg.UseSSL,
		publicURL: cfg.PublicURL,
	}, nil
}

// Upload uploads the file and returns the public URL of the file.
func (u *S3CompatibleUploader) Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	// Automatically check and create bucket if it doesn't exist
	exists, err := u.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = u.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}

		// Set public policy to allow users to directly access files via URL
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = u.client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return "", err
		}
	}

	// Upload file
	_, err = u.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Returns Public URL
	if u.publicURL != "" {
		return u.publicURL + "/" + bucketName + "/" + objectName, nil
	}

	scheme := "http"
	if u.useSSL {
		scheme = "https"
	}

	return scheme + "://" + u.endpoint + "/" + bucketName + "/" + objectName, nil
}
