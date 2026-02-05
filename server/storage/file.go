package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Presigner -> PresignClient, a client used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type Presigner struct {
	PresignClient *s3.PresignClient
}

func NewPresigner(s3Client *s3.Client) *Presigner {
	return &Presigner{
		PresignClient: s3.NewPresignClient(s3Client),
	}
}

func (presigner Presigner) GeneratePresignedURL(ctx context.Context, bucketName string, objectKey string) (string, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	request, err := presigner.PresignClient.PresignGetObject(ctx, params)
	if err != nil {
		return "", err
	}
	return request.URL, err
}
