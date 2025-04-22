package service_test

import (
	"errors"
	"testing"
	"video-processor/internal/mocks"
	"video-processor/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*mocks.BucketClient, *mocks.FrameProcessor, *mocks.FrameRepository, *mocks.NotifierInterface, service.VideoService) {
	mockBucket := new(mocks.BucketClient)
	mockProcessor := new(mocks.FrameProcessor)
	mockRepo := new(mocks.FrameRepository)
	mockNotifier := new(mocks.NotifierInterface)

	svc := service.VideoService{
		Bucket:    mockBucket,
		Processor: mockProcessor,
		Repo:      mockRepo,
		Notifier:  mockNotifier,
	}

	return mockBucket, mockProcessor, mockRepo, mockNotifier, svc
}

func TestProcessVideo_Success(t *testing.T) {
	bucket, processor, repo, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "123",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	bucket.On("DownloadVideo", "video.mp4").Return("/tmp/video.mp4", nil)
	processor.On("ExtractFrames", "/tmp/video.mp4").Return([]string{"frame1.jpg"}, nil)
	bucket.On("UploadFrame", "frame1.jpg", "123").Return("url1", nil)
	repo.On("SaveFrames", "123", "video.mp4", []string{"url1"}).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertNotCalled(t, "NotifyResult", mock.Anything)
}

func TestProcessVideo_ForcedError(t *testing.T) {
	_, _, _, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "forcar-erro",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	notifier.On("NotifyResult", mock.Anything).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertCalled(t, "NotifyResult", mock.Anything)
}

func TestProcessVideo_DownloadError(t *testing.T) {
	bucket, _, _, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "123",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	bucket.On("DownloadVideo", "video.mp4").Return("", errors.New("erro"))
	notifier.On("NotifyResult", mock.Anything).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertCalled(t, "NotifyResult", mock.Anything)
}

func TestProcessVideo_ExtractError(t *testing.T) {
	bucket, processor, _, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "123",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	bucket.On("DownloadVideo", "video.mp4").Return("/tmp/video.mp4", nil)
	processor.On("ExtractFrames", "/tmp/video.mp4").Return(nil, errors.New("erro"))
	notifier.On("NotifyResult", mock.Anything).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertCalled(t, "NotifyResult", mock.Anything)
}

func TestProcessVideo_UploadError(t *testing.T) {
	bucket, processor, _, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "123",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	bucket.On("DownloadVideo", "video.mp4").Return("/tmp/video.mp4", nil)
	processor.On("ExtractFrames", "/tmp/video.mp4").Return([]string{"frame1.jpg"}, nil)
	bucket.On("UploadFrame", "frame1.jpg", "123").Return("", errors.New("erro"))
	notifier.On("NotifyResult", mock.Anything).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertCalled(t, "NotifyResult", mock.Anything)
}

func TestProcessVideo_SaveError(t *testing.T) {
	bucket, processor, repo, notifier, svc := setup()

	event := service.VideoEvent{
		UserID:    "123",
		UserEmail: "user@email.com",
		VideoKey:  "video.mp4",
	}

	bucket.On("DownloadVideo", "video.mp4").Return("/tmp/video.mp4", nil)
	processor.On("ExtractFrames", "/tmp/video.mp4").Return([]string{"frame1.jpg"}, nil)
	bucket.On("UploadFrame", "frame1.jpg", "123").Return("url1", nil)
	repo.On("SaveFrames", "123", "video.mp4", []string{"url1"}).Return(errors.New("erro"))
	notifier.On("NotifyResult", mock.Anything).Return(nil)

	err := svc.ProcessVideo(event)
	assert.NoError(t, err)
	notifier.AssertCalled(t, "NotifyResult", mock.Anything)
}
