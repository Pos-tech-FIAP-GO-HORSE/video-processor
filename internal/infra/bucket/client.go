package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	s3Client *s3.Client
	bucket   = os.Getenv("S3_BUCKET")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
}

// DownloadVideo baixa o vídeo do S3 para o /tmp da Lambda
func DownloadVideo(key string) (string, error) {
	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
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

// UploadFrame faz upload de um frame pro S3 e retorna a URL pública
func UploadFrame(localPath string, userID string) (string, error) {
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

	key := fmt.Sprintf("prints/%s/%s", userID, filepath.Base(localPath))

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String("image/jpeg"),
		ACL:         s3types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return "", err
	}

	// URL pública se o bucket permitir
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
	return url, nil
}
