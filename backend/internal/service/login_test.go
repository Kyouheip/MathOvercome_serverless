package service_test

import (
	"errors"
	"testing"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

// mockLoginRepo は repository.LoginRepo のテスト用実装。
type mockLoginRepo struct {
	findUserByUserIDFn func(userID string) (*model.User, error)
	saveUserFn         func(user *model.User) error
}

func (m *mockLoginRepo) FindUserByUserID(userID string) (*model.User, error) {
	return m.findUserByUserIDFn(userID)
}

func (m *mockLoginRepo) SaveUser(user *model.User) error {
	return m.saveUserFn(user)
}

// --- Authenticate ---

func TestAuthenticate_Success(t *testing.T) {
	repo := &mockLoginRepo{
		findUserByUserIDFn: func(userID string) (*model.User, error) {
			return &model.User{ID: 1, UserID: userID, Password: "pass123"}, nil
		},
	}
	svc := service.NewLoginService(repo)

	user, err := svc.Authenticate(dto.LoginRequest{UserID: "testuser", Password: "pass123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.UserID != "testuser" {
		t.Errorf("expected UserID = testuser, got %s", user.UserID)
	}
}

func TestAuthenticate_UserNotFound(t *testing.T) {
	repo := &mockLoginRepo{
		findUserByUserIDFn: func(userID string) (*model.User, error) {
			return nil, errors.New("record not found")
		},
	}
	svc := service.NewLoginService(repo)

	_, err := svc.Authenticate(dto.LoginRequest{UserID: "unknown", Password: "pass"})
	if !errors.Is(err, apperr.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticate_WrongPassword(t *testing.T) {
	repo := &mockLoginRepo{
		findUserByUserIDFn: func(userID string) (*model.User, error) {
			return &model.User{ID: 1, UserID: userID, Password: "correct"}, nil
		},
	}
	svc := service.NewLoginService(repo)

	_, err := svc.Authenticate(dto.LoginRequest{UserID: "testuser", Password: "wrong"})
	if !errors.Is(err, apperr.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

// --- ValidateRegister ---

func TestValidateRegister_Success(t *testing.T) {
	repo := &mockLoginRepo{
		findUserByUserIDFn: func(userID string) (*model.User, error) {
			return nil, errors.New("record not found")
		},
	}
	svc := service.NewLoginService(repo)

	err := svc.ValidateRegister(dto.RegisterRequest{
		UserID:    "newuser1",
		Password1: "pass123",
		Password2: "pass123",
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateRegister_PasswordMismatch(t *testing.T) {
	repo := &mockLoginRepo{}
	svc := service.NewLoginService(repo)

	err := svc.ValidateRegister(dto.RegisterRequest{
		UserID:    "newuser1",
		Password1: "pass123",
		Password2: "different",
	})
	if !errors.Is(err, apperr.ErrPasswordMismatch) {
		t.Errorf("expected ErrPasswordMismatch, got %v", err)
	}
}

func TestValidateRegister_UserAlreadyExists(t *testing.T) {
	repo := &mockLoginRepo{
		findUserByUserIDFn: func(userID string) (*model.User, error) {
			return &model.User{ID: 1, UserID: userID}, nil
		},
	}
	svc := service.NewLoginService(repo)

	err := svc.ValidateRegister(dto.RegisterRequest{
		UserID:    "existing1",
		Password1: "pass123",
		Password2: "pass123",
	})
	if !errors.Is(err, apperr.ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

// --- CreateUser ---

func TestCreateUser_Success(t *testing.T) {
	var saved *model.User
	repo := &mockLoginRepo{
		saveUserFn: func(user *model.User) error {
			saved = user
			return nil
		},
	}
	svc := service.NewLoginService(repo)

	err := svc.CreateUser(dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "testuser1",
		Password1: "pass123",
		Password2: "pass123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if saved == nil {
		t.Fatal("expected SaveUser to be called")
	}
	if saved.UserID != "testuser1" {
		t.Errorf("expected UserID = testuser1, got %s", saved.UserID)
	}
	if saved.UserName != "Test User" {
		t.Errorf("expected UserName = Test User, got %s", saved.UserName)
	}
	if saved.Password != "pass123" {
		t.Errorf("expected Password = pass123, got %s", saved.Password)
	}
}

func TestCreateUser_RepoError(t *testing.T) {
	repo := &mockLoginRepo{
		saveUserFn: func(user *model.User) error {
			return errors.New("db error")
		},
	}
	svc := service.NewLoginService(repo)

	err := svc.CreateUser(dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "testuser1",
		Password1: "pass123",
		Password2: "pass123",
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
