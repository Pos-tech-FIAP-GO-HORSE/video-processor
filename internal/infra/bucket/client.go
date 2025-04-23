package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3API interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Client struct {
	Client S3API
	Bucket string
}

func (s S3Client) DownloadVideo(key string) (string, error) {
	fmt.Println("Tentando baixar o vídeo:", key, "do bucket:", s.Bucket)

	output, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Erro ao fazer GetObject: %v\n", err)
		return "", err
	}
	defer output.Body.Close()

	filePath := filepath.Join("/tmp", filepath.Base(key))
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, output.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (s S3Client) UploadFrame(localPath string, userID string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("images/%s/%s", userID, filepath.Base(localPath))

	_, err = s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String("image/jpeg"),
		ACL:         s3types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Bucket, key)
	return url, nil
}
