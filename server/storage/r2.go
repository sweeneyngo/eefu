package storage

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	UploadPartSize     = 5 * 1024 * 1024 // 5 MB per part
	UploadConcurrency  = 2               // 2 parts in parallel
	UploadBufferSize   = 5 * 1024 * 1024 // 5 MB buffer for streaming
	UploadCleanupOnErr = true            // clean up parts if upload fails
)

// NOTE: Deprecated: AWS global endpoint resolution interface is deprecated. See deprecation docs on [EndpointResolver].
func NewR2Client(endpoint, accessKey, secretKey string) *s3.Client {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		Region:      "auto", // R2 requires "auto"
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: "auto",
				}, nil
			}),
	}

	return s3.NewFromConfig(cfg)
}

func NewR2Uploader(client *s3.Client) *manager.Uploader {
	return manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = UploadPartSize
		u.Concurrency = UploadConcurrency
		u.LeavePartsOnError = UploadCleanupOnErr
		u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(UploadBufferSize)
	})
}
