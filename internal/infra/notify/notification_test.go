package notify_test

import (
	"context"
	"errors"
	"testing"
	"video-processor/internal/infra/notify"
	"video-processor/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSNSClient struct {
	mock.Mock
}

func (m *MockSNSClient) Publish(ctx context.Context, input *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestNotifyResult_Success(t *testing.T) {
	mockSNS := new(MockSNSClient)

	notifier := notify.Notifier{
		Client:   mockSNS,
		TopicArn: "arn:aws:sns:us-east-1:123456789012:my-topic",
	}

	event := service.NotificationEvent{
		Email:     "user@email.com",
		VideoName: "video.mp4",
	}

	mockSNS.On("Publish", mock.Anything, mock.MatchedBy(func(input *sns.PublishInput) bool {
		return *input.TopicArn == notifier.TopicArn
	})).Return(&sns.PublishOutput{}, nil)

	err := notifier.NotifyResult(event)

	assert.NoError(t, err)
	mockSNS.AssertExpectations(t)
}

func TestNotifyResult_ErrorOnPublish(t *testing.T) {
	mockClient := new(MockSNSClient)
	notifier := notify.Notifier{
		Client:   mockClient,
		TopicArn: "arn:aws:sns:us-east-1:123456789012:test-topic",
	}

	event := service.NotificationEvent{
		Email:     "test@example.com",
		VideoName: "video.mp4",
	}

	mockClient.On("Publish", mock.Anything, mock.Anything).
		Return(&sns.PublishOutput{}, errors.New("falha ao publicar"))

	err := notifier.NotifyResult(event)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falha ao publicar")

	mockClient.AssertExpectations(t)
}
