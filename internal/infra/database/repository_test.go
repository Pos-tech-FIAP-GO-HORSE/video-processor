package database_test

import (
	"errors"
	"testing"
	"video-processor/internal/infra/database"
	"video-processor/internal/mocks"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveFrames_Success(t *testing.T) {
	mockClient := new(mocks.DynamoClient)

	repo := database.Repository{
		Client:    mockClient,
		TableName: "videos_table",
	}

	userID := "123"
	videoKey := "video.mp4"
	urls := []string{"url1", "url2"}

	// Configura o mock para retornar sucesso
	mockClient.
		On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
			return *input.TableName == "videos_table"
		})).
		Return(&dynamodb.PutItemOutput{}, nil)

	err := repo.SaveFrames(userID, videoKey, urls)
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSaveFrames_Error(t *testing.T) {
	mockClient := new(mocks.DynamoClient)

	repo := database.Repository{
		Client:    mockClient,
		TableName: "videos_table",
	}

	mockClient.
		On("PutItem", mock.Anything, mock.Anything).
		Return(nil, errors.New("erro de teste"))

	err := repo.SaveFrames("123", "video.mp4", []string{"url1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao salvar item")
}
