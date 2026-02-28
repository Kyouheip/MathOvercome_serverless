package repository

import "time"

type CategoryStats struct {
	Name         string `gorm:"column:name"`
	TotalCount   int    `gorm:"column:total_count"`
	CorrectCount int    `gorm:"column:correct_count"`
}

func (r *Repository) GetCategoryStats(sessionID uint64) ([]CategoryStats, error) {
	var stats []CategoryStats
	sql := `
		SELECT c.name,
		       COUNT(sp.id) AS total_count,
		       SUM(CASE WHEN sp.is_correct = true THEN 1 ELSE 0 END) AS correct_count
		FROM sessionproblems AS sp
		INNER JOIN problems AS p ON sp.problem_id = p.id
		INNER JOIN categories AS c ON p.category_id = c.id
		WHERE sp.session_id = ?
		GROUP BY c.name`

	err := r.db.Raw(sql, sessionID).Scan(&stats).Error
	return stats, err
}

type SessionProblemRow struct {
	SessionID    uint64    `gorm:"column:session_id"`
	StartTime    time.Time `gorm:"column:start_time"`
	IsCorrect    bool      `gorm:"column:is_correct"`
	CategoryName string    `gorm:"column:category_name"`
}

func (r *Repository) GetSessionProblemsRaw(userID uint64) ([]SessionProblemRow, error) {
	var rows []SessionProblemRow
	sql := `
		SELECT ts.id AS session_id,
		       ts.start_time,
		       sp.is_correct,
		       c.name AS category_name
		FROM test_sessions AS ts
		INNER JOIN sessionproblems AS sp ON sp.session_id = ts.id
		INNER JOIN problems AS p ON sp.problem_id = p.id
		INNER JOIN categories AS c ON p.category_id = c.id
		WHERE ts.user_id = ?
		ORDER BY ts.id DESC, sp.id`

	err := r.db.Raw(sql, userID).Scan(&rows).Error
	return rows, err
}

func (r *Repository) GetWeakCategories(sessionID uint64) ([]string, error) {
	var names []string
	sql := `
		SELECT c.name
		FROM sessionproblems AS sp
		INNER JOIN problems AS p ON sp.problem_id = p.id
		INNER JOIN categories AS c ON p.category_id = c.id
		WHERE sp.session_id = ?
		GROUP BY c.name
		HAVING SUM(CASE WHEN sp.is_correct = true THEN 1 ELSE 0 END) * 1.0 / COUNT(sp.id) < 0.5
		ORDER BY SUM(CASE WHEN sp.is_correct = true THEN 1 ELSE 0 END) * 1.0 / COUNT(sp.id) ASC
		LIMIT 2`

	err := r.db.Raw(sql, sessionID).Scan(&names).Error
	return names, err
}
