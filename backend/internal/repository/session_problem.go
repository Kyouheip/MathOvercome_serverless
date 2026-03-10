package repository

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

type dynamoSP struct {
	PK               string  `dynamodbav:"pk"`
	SK               string  `dynamodbav:"sk"`
	ID               uint64  `dynamodbav:"id"`
	SessionID        uint64  `dynamodbav:"session_id"`
	ProblemID        uint64  `dynamodbav:"problem_id"`
	CategoryID       int     `dynamodbav:"category_id"`
	CategoryName     string  `dynamodbav:"category_name"`
	SelectedChoiceID *uint64 `dynamodbav:"selected_choice_id,omitempty"`
	IsCorrect        *bool   `dynamodbav:"is_correct,omitempty"`
}

// querySessionProblems は pk=SESSION#<id>, sk begins_with SP# で全SPを取得し SK 昇順で返す。
func (r *Repository) querySessionProblems(sessionID uint64) ([]dynamoSP, error) {
	out, err := r.client.Query(bg(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName()),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%d", sessionID)},
			":prefix": &types.AttributeValueMemberS{Value: "SP#"},
		},
	})
	if err != nil {
		return nil, err
	}

	var sps []dynamoSP
	for _, item := range out.Items {
		var dsp dynamoSP
		if err := attributevalue.UnmarshalMap(item, &dsp); err != nil {
			return nil, err
		}
		sps = append(sps, dsp)
	}
	sort.Slice(sps, func(i, j int) bool {
		return sps[i].SK < sps[j].SK
	})
	return sps, nil
}

// fetchProblemWithChoices は pk=PROBLEM#<id> の全アイテム（#METADATA + CHOICE#*）を一度に取得する。
func (r *Repository) fetchProblemWithChoices(problemID uint64) (*model.Problem, []model.Choice, error) {
	out, err := r.client.Query(bg(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName()),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PROBLEM#%d", problemID)},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	var problem *model.Problem
	var choices []model.Choice

	for _, item := range out.Items {
		var skHolder struct {
			SK string `dynamodbav:"sk"`
		}
		if err := attributevalue.UnmarshalMap(item, &skHolder); err != nil {
			return nil, nil, err
		}

		switch {
		case skHolder.SK == "#METADATA":
			var dp dynamoProblem
			if err := attributevalue.UnmarshalMap(item, &dp); err != nil {
				return nil, nil, err
			}
			problem = &model.Problem{
				ID:         dp.ID,
				CategoryID: dp.CategoryID,
				Question:   dp.Question,
				Hint:       dp.Hint,
			}
		case strings.HasPrefix(skHolder.SK, "CHOICE#"):
			var dc dynamoChoice
			if err := attributevalue.UnmarshalMap(item, &dc); err != nil {
				return nil, nil, err
			}
			choices = append(choices, model.Choice{
				ID:         dc.ID,
				ProblemID:  dc.ProblemID,
				ChoiceText: dc.ChoiceText,
				IsCorrect:  dc.IsCorrect,
			})
		}
	}

	if problem == nil {
		return nil, nil, fmt.Errorf("problem not found: %d", problemID)
	}
	return problem, choices, nil
}

func toModelSP(dsp dynamoSP) model.SessionProblem {
	return model.SessionProblem{
		ID:               dsp.ID,
		TestSessionID:    dsp.SessionID,
		ProblemID:        dsp.ProblemID,
		CategoryID:       dsp.CategoryID,
		CategoryName:     dsp.CategoryName,
		SelectedChoiceID: dsp.SelectedChoiceID,
		IsCorrect:        dsp.IsCorrect,
	}
}

func (r *Repository) FindSessionProblemByIdx(sessionID uint64, idx int) (*model.SessionProblem, error) {
	sps, err := r.querySessionProblems(sessionID)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(sps) {
		return nil, fmt.Errorf("index out of range: %d", idx)
	}

	sp := toModelSP(sps[idx])

	problem, choices, err := r.fetchProblemWithChoices(sp.ProblemID)
	if err != nil {
		return nil, err
	}
	problem.Choices = choices
	sp.Problem = *problem

	return &sp, nil
}

func (r *Repository) CountSessionProblems(sessionID uint64) (int64, error) {
	out, err := r.client.Query(bg(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName()),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("SESSION#%d", sessionID)},
			":prefix": &types.AttributeValueMemberS{Value: "SP#"},
		},
		Select: types.SelectCount,
	})
	if err != nil {
		return 0, err
	}
	return int64(out.Count), nil
}

func (r *Repository) FindSessionProblemsBySessionID(sessionID uint64) ([]model.SessionProblem, error) {
	sps, err := r.querySessionProblems(sessionID)
	if err != nil {
		return nil, err
	}
	result := make([]model.SessionProblem, len(sps))
	for i, dsp := range sps {
		result[i] = toModelSP(dsp)
	}
	return result, nil
}

func (r *Repository) SaveSessionProblem(sp *model.SessionProblem) error {
	dsp := dynamoSP{
		PK:               fmt.Sprintf("SESSION#%d", sp.TestSessionID),
		SK:               fmt.Sprintf("SP#%d", sp.ID),
		ID:               sp.ID,
		SessionID:        sp.TestSessionID,
		ProblemID:        sp.ProblemID,
		CategoryID:       sp.CategoryID,
		CategoryName:     sp.CategoryName,
		SelectedChoiceID: sp.SelectedChoiceID,
		IsCorrect:        sp.IsCorrect,
	}
	item, err := attributevalue.MarshalMap(dsp)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(bg(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName()),
		Item:      item,
	})
	return err
}

func (r *Repository) SaveSessionProblems(sps []model.SessionProblem) error {
	if len(sps) == 0 {
		return nil
	}

	// BatchGetItem で問題のカテゴリ情報を一括取得
	keys := make([]map[string]types.AttributeValue, len(sps))
	for i, sp := range sps {
		keys[i] = map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PROBLEM#%d", sp.ProblemID)},
			"sk": &types.AttributeValueMemberS{Value: "#METADATA"},
		}
	}
	batchOut, err := r.client.BatchGetItem(bg(), &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			tableName(): {Keys: keys},
		},
	})
	if err != nil {
		return err
	}

	// problem_id → category_id のマップを構築
	catIDMap := make(map[uint64]int)
	for _, item := range batchOut.Responses[tableName()] {
		var dp dynamoProblem
		if err := attributevalue.UnmarshalMap(item, &dp); err != nil {
			return err
		}
		catIDMap[dp.ID] = dp.CategoryID
	}

	// ID を採番して書き込みリクエストを生成
	base := uint64(time.Now().UnixNano())
	requests := make([]types.WriteRequest, len(sps))
	for i := range sps {
		sps[i].ID = base + uint64(i)
		catID := catIDMap[sps[i].ProblemID]
		dsp := dynamoSP{
			PK:           fmt.Sprintf("SESSION#%d", sps[i].TestSessionID),
			SK:           fmt.Sprintf("SP#%d", sps[i].ID),
			ID:           sps[i].ID,
			SessionID:    sps[i].TestSessionID,
			ProblemID:    sps[i].ProblemID,
			CategoryID:   catID,
			CategoryName: catNames[catID],
		}
		item, err := attributevalue.MarshalMap(dsp)
		if err != nil {
			return err
		}
		requests[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
	}

	// DynamoDB は BatchWriteItem 1回あたり最大 25 件
	for i := 0; i < len(requests); i += 25 {
		end := i + 25
		if end > len(requests) {
			end = len(requests)
		}
		_, err := r.client.BatchWriteItem(bg(), &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName(): requests[i:end],
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
