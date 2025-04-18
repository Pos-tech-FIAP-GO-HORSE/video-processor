package service

import (
	"fmt"
	"log"
	"video-processor/internal/infra/bucket"
	"video-processor/internal/infra/database"
	"video-processor/internal/infra/notify"
	"video-processor/internal/infra/processor"
)

type VideoEvent struct {
	UserID   string `json:"user_id"`
	VideoKey string `json:"video_key"`
}

func ProcessVideo(event VideoEvent) error {
	log.Printf("Processando vídeo: %s", event.VideoKey)

	if event.UserID == "forcar-erro" {
		return notify.NotifyError(notify.VideoEvent(event), fmt.Errorf("erro simulado para teste"))
	}

	videoPath, err := bucket.DownloadVideo(event.VideoKey)
	if err != nil {
		notify.NotifyError(notify.VideoEvent(event), err)
		return err
	}

	frames, err := processor.ExtractFrames(videoPath)
	if err != nil {
		notify.NotifyError(notify.VideoEvent(event), err)
		return err
	}

	var urls []string
	for _, frame := range frames {
		url, err := bucket.UploadFrame(frame, event.UserID)
		if err != nil {
			notify.NotifyError(notify.VideoEvent(event), err)
			return err
		}
		urls = append(urls, url)
	}

	if err := database.SaveFrames(event.UserID, event.VideoKey, urls); err != nil {
		notify.NotifyError(notify.VideoEvent(event), err)
		return err
	}

	if err := notify.NotifySuccess(notify.VideoEvent(event)); err != nil {
		log.Printf("Erro ao notificar sucesso: %v", err)
	}
	return nil
}
