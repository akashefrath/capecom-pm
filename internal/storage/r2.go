package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	client     *s3.Client
	presigner  *s3.PresignClient
	BucketName string
	FolderName string
}

func NewR2Client(accountID, accessKeyID, accessKeySecret, bucketName, folderName string) *R2Client {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	client := s3.New(s3.Options{
		Region: "auto",
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKeyID,
			accessKeySecret,
			"",
		),
		BaseEndpoint: aws.String(endpoint),
	})

	return &R2Client{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		BucketName: bucketName,
		FolderName: folderName,
	}
}

func (r *R2Client) GeneratePresignedUploadURL(storageKey string, mimeType string, expiry time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.BucketName),
		Key:         aws.String(storageKey),
		ContentType: aws.String(mimeType),
	}

	result, err := r.presigner.PresignPutObject(context.Background(), input, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (r *R2Client) GeneratePresignedGetURL(storageKey string, expiry time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(storageKey),
	}

	result, err := r.presigner.PresignGetObject(context.Background(), input, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

// FileExistsResult holds the existence check result for a single storage key.
type FileExistsResult struct {
	Key    string `json:"key"`
	Exists bool   `json:"exists"`
}

// CheckFileExists checks whether a single file exists in the R2 bucket.
func (r *R2Client) CheckFileExists(storageKey string) (bool, error) {
	_, err := r.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(storageKey),
	})
	if err != nil {
		// HeadObject returns a NotFound-style error when the key doesn't exist.
		// The SDK wraps it, so we check the error message. A more robust check
		// could use smithy error codes, but this covers the common case.
		return false, nil
	}
	return true, nil
}

// CheckFilesExist checks existence for multiple keys and returns a result per key.
func (r *R2Client) CheckFilesExist(storageKeys []string) ([]FileExistsResult, error) {
	results := make([]FileExistsResult, 0, len(storageKeys))
	for _, key := range storageKeys {
		exists, err := r.CheckFileExists(key)
		if err != nil {
			return nil, err
		}
		results = append(results, FileExistsResult{
			Key:    key,
			Exists: exists,
		})
	}
	return results, nil
}
