package handlers_test

import (
	"encoding/json"
	"errors"
	"testing"
	"video-processor/internal/handlers"
	"video-processor/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	CalledWith service.VideoEvent
	Err        error
}

func (m *mockService) ProcessVideo(event service.VideoEvent) error {
	m.CalledWith = event
	return m.Err
}

func TestHandleRequest_Success(t *testing.T) {
	mockSvc := &mockService{}
	handler := handlers.NewHandler(mockSvc)

	input := handlers.VideoEvent{
		UserID:    "123",
		UserEmail: "user@example.com",
		VideoKey:  "video.mp4",
	}
	msg, _ := json.Marshal(input)

	snsEvent := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: string(msg),
				},
			},
		},
	}

	err := handler.HandleRequest(snsEvent)
	assert.NoError(t, err)
	assert.Equal(t, mockSvc.CalledWith.UserID, "123")
}

func TestHandleRequest_ParseError(t *testing.T) {
	mockSvc := &mockService{}
	handler := handlers.NewHandler(mockSvc)

	snsEvent := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: "{invalid-json",
				},
			},
		},
	}

	err := handler.HandleRequest(snsEvent)
	assert.Error(t, err)
}

func TestHandleRequest_ProcessVideoError(t *testing.T) {
	mockSvc := &mockService{Err: errors.New("erro no processamento")}
	handler := handlers.NewHandler(mockSvc)

	input := handlers.VideoEvent{
		UserID:    "456",
		UserEmail: "erro@example.com",
		VideoKey:  "fail.mp4",
	}
	msg, _ := json.Marshal(input)

	snsEvent := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message: string(msg),
				},
			},
		},
	}

	err := handler.HandleRequest(snsEvent)
	assert.NoError(t, err)
	assert.Equal(t, mockSvc.CalledWith.VideoKey, "fail.mp4")
}
