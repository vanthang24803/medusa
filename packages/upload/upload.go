package upload

import (
	"context"
	"fmt"
	"io"
	"mime"
	"strings"

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
	BucketName      string // Default bucket name; overrides the bucketName param in Upload when set
}

// S3CompatibleUploader implements Uploader interface for all S3-compatible storage.
type S3CompatibleUploader struct {
	client    *minio.Client
	provider  ProviderType
	endpoint  string
	useSSL    bool
	publicURL string
	bucket    string // default bucket; used when Upload is called with an empty bucketName
}

// nopUploader is a no-op uploader returned when the configured provider is unavailable.
type nopUploader struct{}

func NewNopUploader() Uploader { return &nopUploader{} }

func (*nopUploader) Upload(_ context.Context, _, _ string, _ io.Reader, _ int64, _ string) (string, error) {
	return "", fmt.Errorf("upload provider not configured")
}

// ExtFromContentType returns a file extension for the given MIME type.
func ExtFromContentType(ct string) string {
	mt, _, _ := mime.ParseMediaType(ct)
	switch strings.ToLower(mt) {
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".jpg"
	}
}

// ValidateSize returns an error if size exceeds maxBytes.
func ValidateSize(size, maxBytes int64) error {
	if size > maxBytes {
		return fmt.Errorf("file size %d bytes exceeds limit of %d bytes", size, maxBytes)
	}
	return nil
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
		provider:  cfg.Provider,
		endpoint:  cfg.Endpoint,
		useSSL:    cfg.UseSSL,
		publicURL: cfg.PublicURL,
		bucket:    cfg.BucketName,
	}, nil
}

// Upload uploads the file and returns the public URL of the file.
// If the uploader was configured with a BucketName, it takes precedence over the bucketName argument.
func (u *S3CompatibleUploader) Upload(ctx context.Context, bucketName string, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	if u.bucket != "" {
		bucketName = u.bucket
	}

	// Check if bucket exists; if not, attempt to create it (best-effort — providers like R2 require
	// buckets to be pre-created via their dashboard, so creation errors are non-fatal here).
	exists, err := u.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		if err = u.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err == nil {
			policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
			_ = u.client.SetBucketPolicy(ctx, bucketName, policy)
		}
		// If MakeBucket failed the bucket may already exist on the provider side; proceed with upload.
	}

	_, err = u.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	if u.publicURL != "" {
		// R2 public bucket URLs are bucket-specific; the bucket name is not part of the object path.
		// For MinIO/S3, the public URL is a base URL that requires the bucket prefix.
		if u.provider == ProviderR2 {
			return u.publicURL + "/" + objectName, nil
		}
		return u.publicURL + "/" + bucketName + "/" + objectName, nil
	}

	scheme := "http"
	if u.useSSL {
		scheme = "https"
	}
	return scheme + "://" + u.endpoint + "/" + bucketName + "/" + objectName, nil
}
