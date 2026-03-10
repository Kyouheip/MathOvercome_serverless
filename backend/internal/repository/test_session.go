package repository

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

type dynamoSession struct {
	PK              string `dynamodbav:"pk"`
	SK              string `dynamodbav:"sk"`
	GSI1PK          string `dynamodbav:"gsi1pk"`
	GSI1SK          string `dynamodbav:"gsi1sk"`
	ID              uint64 `dynamodbav:"id"`
	UserID          uint64 `dynamodbav:"user_id"`
	IncludeIntegers bool   `dynamodbav:"include_integers"`
	StartTime       string `dynamodbav:"start_time"`
}

func (r *Repository) FindTestSessionByID(id uint64) (*model.TestSession, error) {
	out, err := r.client.GetItem(bg(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName()),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%d", id)},
			"sk": &types.AttributeValueMemberS{Value: "#METADATA"},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("session not found: %d", id)
	}

	var ds dynamoSession
	if err := attributevalue.UnmarshalMap(out.Item, &ds); err != nil {
		return nil, err
	}

	startTime, _ := time.Parse("2006-01-02 15:04:05", ds.StartTime)
	return &model.TestSession{
		ID:              ds.ID,
		UserID:          ds.UserID,
		IncludeIntegers: ds.IncludeIntegers,
		StartTime:       startTime,
	}, nil
}

func (r *Repository) SaveTestSession(session *model.TestSession) error {
	session.ID = uint64(time.Now().UnixNano())
	session.StartTime = time.Now()

	ds := dynamoSession{
		PK:              fmt.Sprintf("SESSION#%d", session.ID),
		SK:              "#METADATA",
		GSI1PK:          fmt.Sprintf("USER#%d", session.UserID),
		GSI1SK:          fmt.Sprintf("SESSION#%d", session.ID),
		ID:              session.ID,
		UserID:          session.UserID,
		IncludeIntegers: session.IncludeIntegers,
		StartTime:       session.StartTime.Format("2006-01-02 15:04:05"),
	}
	item, err := attributevalue.MarshalMap(ds)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(bg(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName()),
		Item:      item,
	})
	return err
}
