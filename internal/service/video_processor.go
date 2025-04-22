package service

import (
	"log"
)

type VideoEvent struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	VideoKey  string `json:"video_key"`
}

type BucketClient interface {
	DownloadVideo(key string) (string, error)
	UploadFrame(localPath string, userID string) (string, error)
}

type FrameProcessor interface {
	ExtractFrames(videoPath string) ([]string, error)
}

type FrameRepository interface {
	SaveFrames(userID, videoKey string, urls []string) error
}

type Notifier interface {
	NotifyResult(event NotificationEvent) error
}

type NotificationEvent struct {
	Email     string `json:"email"`
	VideoName string `json:"videoName"`
}

type VideoService struct {
	Bucket    BucketClient
	Processor FrameProcessor
	Repo      FrameRepository
	Notifier  Notifier
}

func (vs VideoService) ProcessVideo(event VideoEvent) error {
	log.Printf("Processando vídeo: %s", event.VideoKey)
	log.Printf("Teste: %s", event.VideoKey)

	notifyEvent := NotificationEvent{
		Email:     event.UserEmail,
		VideoName: event.VideoKey,
	}

	if event.UserID == "forcar-erro" {
		log.Printf("Vai notificar")
		return vs.Notifier.NotifyResult(notifyEvent)
	}

	log.Printf("Vai fazer o download")
	videoPath, err := vs.Bucket.DownloadVideo(event.VideoKey)
	if err != nil {
		return vs.Notifier.NotifyResult(notifyEvent)
	}

	frames, err := vs.Processor.ExtractFrames(videoPath)
	if err != nil {
		return vs.Notifier.NotifyResult(notifyEvent)
	}

	var urls []string
	for _, frame := range frames {
		url, err := vs.Bucket.UploadFrame(frame, event.UserID)
		if err != nil {
			return vs.Notifier.NotifyResult(notifyEvent)
		}
		urls = append(urls, url)
	}

	if err := vs.Repo.SaveFrames(event.UserID, event.VideoKey, urls); err != nil {
		return vs.Notifier.NotifyResult(notifyEvent)
	}

	return nil
}
