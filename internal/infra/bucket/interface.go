package bucket

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketClient interface {
	DownloadVideo(key string) (string, error)
	UploadFrame(localPath string, userID string) (string, error)
}

type BucketS3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}
