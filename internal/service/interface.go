package service

type VideoProcessor interface {
	ProcessVideo(event VideoEvent) error
}
