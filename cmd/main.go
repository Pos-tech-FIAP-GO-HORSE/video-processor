package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"log"
	"os"
	"video-processor/internal/handlers"
	"video-processor/internal/infra/bucket"
	"video-processor/internal/infra/database"
	"video-processor/internal/infra/notify"
	"video-processor/internal/infra/processor"
	"video-processor/internal/service"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS: %v", err)
	}

	notifier := notify.Notifier{
		Client:   sns.NewFromConfig(cfg),
		TopicArn: os.Getenv("USER_NOTIFICATION_TOPIC_ARN"),
	}

	s3Client := &bucket.S3Client{
		Client: s3.NewFromConfig(cfg),
		Bucket: os.Getenv("S3_BUCKET"),
	}

	repo := database.Repository{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: os.Getenv("DYNAMODB_TABLE"),
	}

	videoService := service.VideoService{
		Bucket:    s3Client,
		Processor: processor.FFmpegProcessor{},
		Repo:      repo,
		Notifier:  notifier,
	}

	handler := handlers.NewHandler(&videoService)
	lambda.Start(handler.HandleRequest)
}
