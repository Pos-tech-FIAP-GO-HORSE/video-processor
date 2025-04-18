package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	client    *dynamodb.Client
	tableName = os.Getenv("DYNAMODB_TABLE")
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Erro ao carregar config AWS: %v", err)
	}
	client = dynamodb.NewFromConfig(cfg)
}

func SaveFrames(userID, videoKey string, urls []string) error {
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

	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("erro ao salvar item no DynamoDB: %v", err)
	}

	return nil
}
