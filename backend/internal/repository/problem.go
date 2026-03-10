package repository

import (
	"fmt"
	"math/rand"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

type dynamoProblem struct {
	PK         string `dynamodbav:"pk"`
	SK         string `dynamodbav:"sk"`
	ID         uint64 `dynamodbav:"id"`
	CategoryID int    `dynamodbav:"category_id"`
	Question   string `dynamodbav:"question"`
	Hint       string `dynamodbav:"hint"`
}

// FindProblemsPerCategory は GSI1 でカテゴリ別に問題を取得し、
// カテゴリごとに countPerCategory 件をランダムに返す。
func (r *Repository) FindProblemsPerCategory(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
	var result []model.Problem

	for _, catID := range categoryIDs {
		out, err := r.client.Query(bg(), &dynamodb.QueryInput{
			TableName:              aws.String(tableName()),
			IndexName:              aws.String("GSI1"),
			KeyConditionExpression: aws.String("gsi1pk = :gsi1pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":gsi1pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CATEGORY#%d", catID)},
			},
		})
		if err != nil {
			return nil, err
		}

		var problems []model.Problem
		for _, item := range out.Items {
			var dp dynamoProblem
			if err := attributevalue.UnmarshalMap(item, &dp); err != nil {
				return nil, err
			}
			problems = append(problems, model.Problem{
				ID:         dp.ID,
				CategoryID: dp.CategoryID,
				Question:   dp.Question,
				Hint:       dp.Hint,
			})
		}

		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})

		take := countPerCategory
		if take > len(problems) {
			take = len(problems)
		}
		result = append(result, problems[:take]...)
	}

	return result, nil
}
