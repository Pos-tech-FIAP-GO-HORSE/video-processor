package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	snsClient             *sns.Client
	userNotificationTopic = os.Getenv("USER_NOTIFICATION_TOPIC_ARN")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS para SNS: %v", err)
	}
	snsClient = sns.NewFromConfig(cfg)
}

type VideoEvent struct {
	Email     string `json:"email"`
	VideoName string `json:"videoName"`
}

func NotifyResult(event VideoEvent) error {
	body, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		return fmt.Errorf("erro ao serializar evento para SNS: %w", marshalErr)
	}

	_, pubErr := snsClient.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(userNotificationTopic),
		Message:  aws.String(string(body)),
	})
	if pubErr != nil {
		log.Printf("Erro ao publicar notificação no SNS: %v", pubErr)
	}
	return pubErr
}
