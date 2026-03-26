package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

// TestSessionRepo は TestSessionService が使うリポジトリ操作を定義する。
type TestSessionRepo interface {
	SaveTestSession(session *model.TestSession) error
	FindProblemsPerCategory(categoryIDs []int, countPerCategory int) ([]model.Problem, error)
	SaveSessionProblems(sps []model.SessionProblem) error
	CountSessionProblems(sessionID uint64) (int64, error)
	FindSessionProblemByIdx(sessionID uint64, idx int) (*model.SessionProblem, error)
	FindSessionProblemsBySessionID(sessionID uint64) ([]model.SessionProblem, error)
	FindChoiceByProblemAndChoiceID(problemID, choiceID uint64) (*model.Choice, error)
	SaveSessionProblem(sp *model.SessionProblem) error
}

// MypageRepo は MypageService が使うリポジトリ操作を定義する。
type MypageRepo interface {
	GetSessionProblemsRaw(userSub string) ([]SessionProblemRow, error)
	GetCategoryStats(sessionID uint64) ([]CategoryStats, error)
	GetWeakCategories(sessionID uint64) ([]string, error)
}
