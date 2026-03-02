package repository

import (
	"fmt"
	"strings"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

func (r *Repository) FindProblemsPerCategory(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
	var queries []string
	for _, id := range categoryIDs {
		queries = append(queries, fmt.Sprintf(
			"(SELECT id FROM problems WHERE category_id = %d ORDER BY RAND() LIMIT %d)",
			id, countPerCategory,
		))
	}

	fullQuery := strings.Join(queries, " UNION ALL ")

	var problemIDs []uint64
	if err := r.db.Raw(fullQuery).Scan(&problemIDs).Error; err != nil {
		return nil, err
	}

	var problems []model.Problem
	err := r.db.Where("id IN ?", problemIDs).Find(&problems).Error
	return problems, err
}
