package service

import (
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
)

// LoginServicer は認証操作を定義する。
type LoginServicer interface {
	Authenticate(req dto.LoginRequest) (*model.User, error)
	ValidateRegister(req dto.RegisterRequest) error
	CreateUser(req dto.RegisterRequest) error
}

// TestSessionServicer はテストセッション操作を定義する。
type TestSessionServicer interface {
	CreateTestSess(user *model.User, includeIntegers bool) (*model.TestSession, error)
}

// MypageServicer はマイページ操作を定義する。
type MypageServicer interface {
	GetUserData(user *model.User) (*dto.User, error)
}
