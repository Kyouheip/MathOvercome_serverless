package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

type dynamoChoice struct {
	PK         string `dynamodbav:"pk"`
	SK         string `dynamodbav:"sk"`
	ID         uint64 `dynamodbav:"id"`
	ProblemID  uint64 `dynamodbav:"problem_id"`
	ChoiceText string `dynamodbav:"choice_text"`
	IsCorrect  bool   `dynamodbav:"is_correct"`
}

// FindChoiceByProblemAndChoiceID は pk=PROBLEM#<problemID>, sk=CHOICE#<choiceID> で選択肢を取得する。
func (r *Repository) FindChoiceByProblemAndChoiceID(problemID, choiceID uint64) (*model.Choice, error) {
	out, err := r.client.GetItem(bg(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PROBLEM#%d", problemID)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("CHOICE#%d", choiceID)},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("choice not found: problem=%d choice=%d", problemID, choiceID)
	}

	var dc dynamoChoice
	if err := attributevalue.UnmarshalMap(out.Item, &dc); err != nil {
		return nil, err
	}
	return &model.Choice{
		ID:         dc.ID,
		ProblemID:  dc.ProblemID,
		ChoiceText: dc.ChoiceText,
		IsCorrect:  dc.IsCorrect,
	}, nil
}
