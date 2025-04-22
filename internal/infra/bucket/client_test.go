package bucket_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"testing"
	"video-processor/internal/infra/bucket"
	"video-processor/internal/mocks"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestDownloadVideo_Success(t *testing.T) {
	mockS3 := new(mocks.S3API)

	mockS3.On("GetObject", mock.Anything, mock.Anything).
		Return(&s3.GetObjectOutput{
			Body: io.NopCloser(bytes.NewReader([]byte("fake content"))),
		}, nil)

	client := bucket.S3Client{
		Client: mockS3,
		Bucket: "my-bucket",
	}

	filePath, err := client.DownloadVideo("video.mp4")

	assert.NoError(t, err)
	assert.Contains(t, filePath, "/tmp/video.mp4")

	mockS3.AssertExpectations(t)
}

func TestDownloadVideo_Error(t *testing.T) {
	mockS3 := new(mocks.S3API)

	mockS3.On("GetObject", mock.Anything, mock.Anything).
		Return(nil, errors.New("falha no download"))

	client := bucket.S3Client{
		Client: mockS3,
		Bucket: "my-bucket",
	}

	_, err := client.DownloadVideo("video.mp4")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falha no download")
}

func TestUploadFrame_Success(t *testing.T) {
	mockS3 := new(mocks.S3API)

	mockS3.On("PutObject", mock.Anything, mock.Anything).
		Return(&s3.PutObjectOutput{}, nil)

	client := bucket.S3Client{
		Client: mockS3,
		Bucket: "my-bucket",
	}

	tmpFile, err := os.CreateTemp("/tmp", "frame-*.jpg")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("fake image data"))
	assert.NoError(t, err)
	tmpFile.Close()

	url, err := client.UploadFrame(tmpFile.Name(), "123")
	assert.NoError(t, err)
	assert.Contains(t, url, "https://my-bucket.s3.amazonaws.com/images/123")
}

func TestUploadFrame_Error(t *testing.T) {
	mockS3 := new(mocks.S3API)

	mockS3.On("PutObject", mock.Anything, mock.Anything).
		Return(nil, errors.New("erro no upload"))

	client := bucket.S3Client{
		Client: mockS3,
		Bucket: "my-bucket",
	}

	tmpFile, err := os.CreateTemp("/tmp", "frame-*.jpg")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = client.UploadFrame(tmpFile.Name(), "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro no upload")
}
