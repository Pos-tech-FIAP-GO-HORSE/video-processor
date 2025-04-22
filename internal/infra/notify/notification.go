package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"video-processor/internal/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSClient interface {
	Publish(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type Notifier struct {
	Client   SNSClient
	TopicArn string
}

func (n Notifier) NotifyResult(event service.NotificationEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento para SNS: %w", err)
	}

	_, pubErr := n.Client.Publish(context.TODO(), &sns.PublishInput{
		TopicArn: aws.String(n.TopicArn),
		Message:  aws.String(string(body)),
	})
	if pubErr != nil {
		log.Printf("Erro ao publicar notificação no SNS: %v", pubErr)
	}
	return pubErr
}
