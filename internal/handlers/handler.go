package handlers

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"video-processor/internal/service"
)

type VideoEvent struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	VideoKey  string `json:"video_key"`
}

func HandleRequest(snsEvent events.SNSEvent) error {
	for _, record := range snsEvent.Records {
		var videoEvent VideoEvent
		err := json.Unmarshal([]byte(record.SNS.Message), &videoEvent)
		if err != nil {
			log.Printf("Erro ao fazer parse do evento SNS: %v", err)
			return err
		}

		if err := service.ProcessVideo(service.VideoEvent(videoEvent)); err != nil {
			log.Printf("Erro ao processar vídeo: %v", err)
		}
	}
	return nil
}
