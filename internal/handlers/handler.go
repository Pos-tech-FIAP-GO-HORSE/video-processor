package handlers

import (
	"context"
	"encoding/json"
	"log"
	"video-processor/internal/infra/bucket"
	"video-processor/internal/infra/database/dynamodb"
	"video-processor/internal/infra/notify"
	"video-processor/internal/infra/processor"

	"github.com/aws/aws-lambda-go/events"
)

type VideoEvent struct {
	UserID   string `json:"userId"`
	VideoKey string `json:"videoKey"`
}

func HandleRequest(ctx context.Context, snsEvent events.SNSEvent) error {
	for _, record := range snsEvent.Records {
		var videoEvent VideoEvent
		err := json.Unmarshal([]byte(record.SNS.Message), &videoEvent)
		if err != nil {
			log.Printf("Erro ao fazer parse do evento SNS: %v", err)
			continue
		}

		log.Printf("Processando vídeo: %s", videoEvent.VideoKey)

		// 1. Baixa o vídeo do S3
		videoPath, err := bucket.DownloadVideo(videoEvent.VideoKey)
		if err != nil {
			log.Printf("Erro ao baixar vídeo: %v", err)
			notify.NotifyError(videoEvent, err)
			continue
		}

		// 2. Extrai imagens com ffmpeg
		frames, err := processor.ExtractFrames(videoPath)
		if err != nil {
			log.Printf("Erro ao extrair frames: %v", err)
			notify.NotifyError(videoEvent, err)
			continue
		}

		// 3. Faz upload das imagens no S3 e salva URLs
		var urls []string
		for _, frame := range frames {
			url, err := s3.UploadFrame(frame, videoEvent.UserID)
			if err != nil {
				log.Printf("Erro ao subir frame: %v", err)
				continue
			}
			urls = append(urls, url)
		}

		// 4. Salva URLs no banco
		err = dynamodb.SaveFrames(videoEvent.UserID, videoEvent.VideoKey, urls)
		if err != nil {
			log.Printf("Erro ao salvar no banco: %v", err)
			continue
		}

		// 5. Notifica sucesso
		err = notify.NotifySuccess(videoEvent)
		if err != nil {
			log.Printf("Erro ao notificar sucesso: %v", err)
		}
	}
	return nil
}
