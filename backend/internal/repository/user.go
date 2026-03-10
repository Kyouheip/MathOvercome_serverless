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

type dynamoUser struct {
	PK       string `dynamodbav:"pk"`
	SK       string `dynamodbav:"sk"`
	GSI1PK   string `dynamodbav:"gsi1pk"`
	GSI1SK   string `dynamodbav:"gsi1sk"`
	ID       uint64 `dynamodbav:"id"`
	UserName string `dynamodbav:"user_name"`
	UserID   string `dynamodbav:"user_id"`
	Password string `dynamodbav:"password"`
}

// FindUserByUserID は GSI2 (user_id = ログイン文字列) でユーザーを検索する。
func (r *Repository) FindUserByUserID(userID string) (*model.User, error) {
	out, err := r.client.Query(bg(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName()),
		IndexName:              aws.String("GSI2"),
		KeyConditionExpression: aws.String("user_id = :uid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uid": &types.AttributeValueMemberS{Value: userID},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	var du dynamoUser
	if err := attributevalue.UnmarshalMap(out.Items[0], &du); err != nil {
		return nil, err
	}
	return &model.User{
		ID:       du.ID,
		UserName: du.UserName,
		UserID:   du.UserID,
		Password: du.Password,
	}, nil
}

func (r *Repository) SaveUser(user *model.User) error {
	user.ID = uint64(time.Now().UnixNano())
	du := dynamoUser{
		PK:       fmt.Sprintf("USER#%d", user.ID),
		SK:       "#METADATA",
		GSI1PK:   "USER",
		GSI1SK:   fmt.Sprintf("USER#%d", user.ID),
		ID:       user.ID,
		UserName: user.UserName,
		UserID:   user.UserID,
		Password: user.Password,
	}
	item, err := attributevalue.MarshalMap(du)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(bg(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName()),
		Item:      item,
	})
	return err
}
