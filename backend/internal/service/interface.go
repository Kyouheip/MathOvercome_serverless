package service

import (
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

// TestSessionServicer はテストセッション操作を定義する。
type TestSessionServicer interface {
	CreateTestSess(userSub string, includeIntegers bool) (*model.TestSession, error)
	GetProblem(sessionID uint64, idx int) (*dto.SessionProblem, error)
	SubmitAnswer(sessionID uint64, idx int, choiceID *int64) error
}

// MypageServicer はマイページ操作を定義する。
type MypageServicer interface {
	GetUserData(user *model.User) (*dto.User, error)
}
