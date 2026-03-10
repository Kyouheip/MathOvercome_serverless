package repository

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func tableName() string {
	if t := os.Getenv("DYNAMODB_TABLE"); t != "" {
		return t
	}
	return "MathOvercome"
}

type Repository struct {
	client *dynamodb.Client
}

func NewRepository(client *dynamodb.Client) *Repository {
	return &Repository{client: client}
}

func bg() context.Context {
	return context.Background()
}

// カテゴリID→カテゴリ名の固定マッピング
var catNames = map[int]string{
	1: "数と式",
	2: "2次関数",
	3: "図形と計量",
	4: "データの分析",
	5: "確率",
	6: "図形の性質",
	7: "整数",
}
