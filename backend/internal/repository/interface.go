package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

// LoginRepo は LoginService が使うリポジトリ操作を定義する。
type LoginRepo interface {
	FindUserByUserID(userID string) (*model.User, error)
	SaveUser(user *model.User) error
}

// TestSessionRepo は TestSessionService が使うリポジトリ操作を定義する。
type TestSessionRepo interface {
	SaveTestSession(session *model.TestSession) error
	FindProblemsPerCategory(categoryIDs []int, countPerCategory int) ([]model.Problem, error)
	SaveSessionProblems(sps []model.SessionProblem) error
}

// MypageRepo は MypageService が使うリポジトリ操作を定義する。
type MypageRepo interface {
	GetSessionProblemsRaw(userID uint64) ([]SessionProblemRow, error)
	GetCategoryStats(sessionID uint64) ([]CategoryStats, error)
	GetWeakCategories(sessionID uint64) ([]string, error)
}

// SessionRepo は SessionHandler が直接使うリポジトリ操作を定義する。
type SessionRepo interface {
	FindTestSessionByID(id uint64) (*model.TestSession, error)
	CountSessionProblems(sessionID uint64) (int64, error)
	FindSessionProblemByIdx(sessionID uint64, idx int) (*model.SessionProblem, error)
	FindSessionProblemsBySessionID(sessionID uint64) ([]model.SessionProblem, error)
	FindChoiceByID(id uint64) (*model.Choice, error)
	SaveSessionProblem(sp *model.SessionProblem) error
}
