package service

import (
	"log"
	"video-processor/internal/infra/bucket"
	"video-processor/internal/infra/database"
	"video-processor/internal/infra/notify"
	"video-processor/internal/infra/processor"
)

type VideoEvent struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	VideoKey  string `json:"video_key"`
}

func ProcessVideo(event VideoEvent) error {
	log.Printf("Processando vídeo: %s", event.VideoKey)

	notifyEvent := notify.VideoEvent{
		Email:     event.UserEmail,
		VideoName: event.VideoKey,
	}
	//somente para gerar o vídeo
	if event.UserID == "forcar-erro" {
		return notify.NotifyResult(notifyEvent)
	}

	videoPath, err := bucket.DownloadVideo(event.VideoKey)
	if err != nil {
		notify.NotifyResult(notifyEvent)
		return err
	}

	frames, err := processor.ExtractFrames(videoPath)
	if err != nil {
		notify.NotifyResult(notifyEvent)
		return err
	}

	var urls []string
	for _, frame := range frames {
		url, err := bucket.UploadFrame(frame, event.UserID)
		if err != nil {
			notify.NotifyResult(notifyEvent)
			return err
		}
		urls = append(urls, url)
	}

	if err := database.SaveFrames(event.UserID, event.VideoKey, urls); err != nil {
		notify.NotifyResult(notifyEvent)
		return err
	}

	return nil
}
