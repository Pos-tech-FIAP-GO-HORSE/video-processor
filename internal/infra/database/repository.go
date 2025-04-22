package database

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	Client    DynamoClient
	TableName string
}

func (r Repository) SaveFrames(userID, videoKey string, urls []string) error {
	item := map[string]types.AttributeValue{
		"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
		"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("VIDEO#%s", videoKey)},
		"Status":    &types.AttributeValueMemberS{Value: "success"},
		"CreatedAt": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
	}

	var prints []types.AttributeValue
	for _, url := range urls {
		prints = append(prints, &types.AttributeValueMemberS{Value: url})
	}
	item["Prints"] = &types.AttributeValueMemberL{Value: prints}

	_, err := r.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("erro ao salvar item no DynamoDB: %v", err)
	}

	return nil
}
