package repository

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CategoryStats struct {
	Name         string
	TotalCount   int
	CorrectCount int
}

type SessionProblemRow struct {
	SessionID    uint64
	StartTime    time.Time
	IsCorrect    bool
	CategoryName string
}

// GetSessionProblemsRaw はユーザーの全セッション×全SPを結合して返す。
// セッションは降順（新しい順）、SP は昇順。
func (r *Repository) GetSessionProblemsRaw(userSub string) ([]SessionProblemRow, error) {
	// GSI1: gsi1pk = USER#<sub> → セッション一覧を降順で取得
	sessOut, err := r.client.Query(bg(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName()),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("gsi1pk = :gsi1pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi1pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userSub)},
		},
		ScanIndexForward: aws.Bool(false), // 降順
	})
	if err != nil {
		return nil, err
	}

	var rows []SessionProblemRow

	for _, sessItem := range sessOut.Items {
		var ds dynamoSession
		if err := attributevalue.UnmarshalMap(sessItem, &ds); err != nil {
			return nil, err
		}
		startTime, _ := time.Parse("2006-01-02 15:04:05", ds.StartTime)

		sps, err := r.querySessionProblems(ds.ID)
		if err != nil {
			return nil, err
		}

		for _, dsp := range sps {
			isCorrect := false
			if dsp.IsCorrect != nil {
				isCorrect = *dsp.IsCorrect
			}
			rows = append(rows, SessionProblemRow{
				SessionID:    ds.ID,
				StartTime:    startTime,
				IsCorrect:    isCorrect,
				CategoryName: dsp.CategoryName,
			})
		}
	}

	return rows, nil
}

// GetCategoryStats はセッション内のSPをカテゴリ別に集計して返す。
func (r *Repository) GetCategoryStats(sessionID uint64) ([]CategoryStats, error) {
	sps, err := r.querySessionProblems(sessionID)
	if err != nil {
		return nil, err
	}

	statsMap := make(map[string]*CategoryStats)
	var order []string

	for _, dsp := range sps {
		if _, exists := statsMap[dsp.CategoryName]; !exists {
			statsMap[dsp.CategoryName] = &CategoryStats{Name: dsp.CategoryName}
			order = append(order, dsp.CategoryName)
		}
		statsMap[dsp.CategoryName].TotalCount++
		if dsp.IsCorrect != nil && *dsp.IsCorrect {
			statsMap[dsp.CategoryName].CorrectCount++
		}
	}

	result := make([]CategoryStats, len(order))
	for i, name := range order {
		result[i] = *statsMap[name]
	}
	return result, nil
}

// GetWeakCategories は正答率 < 0.5 のカテゴリを昇順で最大2件返す。
func (r *Repository) GetWeakCategories(sessionID uint64) ([]string, error) {
	sps, err := r.querySessionProblems(sessionID)
	if err != nil {
		return nil, err
	}

	type catStat struct {
		total   int
		correct int
	}
	statsMap := make(map[string]*catStat)
	var order []string

	for _, dsp := range sps {
		if _, exists := statsMap[dsp.CategoryName]; !exists {
			statsMap[dsp.CategoryName] = &catStat{}
			order = append(order, dsp.CategoryName)
		}
		statsMap[dsp.CategoryName].total++
		if dsp.IsCorrect != nil && *dsp.IsCorrect {
			statsMap[dsp.CategoryName].correct++
		}
	}

	type weakCat struct {
		name string
		rate float64
	}
	var weaks []weakCat
	for _, name := range order {
		st := statsMap[name]
		if st.total == 0 {
			continue
		}
		rate := float64(st.correct) / float64(st.total)
		if rate < 0.5 {
			weaks = append(weaks, weakCat{name: name, rate: rate})
		}
	}

	sort.Slice(weaks, func(i, j int) bool {
		return weaks[i].rate < weaks[j].rate
	})

	names := []string{}
	for i := 0; i < len(weaks) && i < 2; i++ {
		names = append(names, weaks[i].name)
	}
	return names, nil
}
