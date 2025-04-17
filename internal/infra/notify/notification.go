package notify

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	snsClient       *sns.Client
	successTopicArn = os.Getenv("PROCESSAMENTO_SUCESSO_TOPIC_ARN")
	errorTopicArn   = os.Getenv("PROCESSAMENTO_ERRO_TOPIC_ARN")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS para SNS: %v", err)
	}
	snsClient = sns.NewFromConfig(cfg)
}

type VideoEvent struct {
	UserID   string `json:"userId"`
	VideoKey string `json:"videoKey"`
}

func NotifySuccess(event VideoEvent) error {
	return publish(event, successTopicArn)
}

func NotifyError(event VideoEvent, err error) error {
	payload := map[string]interface{}{
		"userId":   event.UserID,
		"videoKey": event.VideoKey,
		"error":    err.Error(),
	}
	body, _ := json.Marshal(payload)

	_, pubErr := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(errorTopicArn),
		Message:  aws.String(string(body)),
	})
	if pubErr != nil {
		log.Printf("Erro ao notificar erro no SNS: %v", pubErr)
	}
	return pubErr
}

func publish(event VideoEvent, topicArn string) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, pubErr := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(body)),
	})
	if pubErr != nil {
		log.Printf("Erro ao publicar mensagem SNS: %v", pubErr)
	}
	return pubErr
}
