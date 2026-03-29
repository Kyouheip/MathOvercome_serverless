package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

type mockMypageRepo struct {
	getSessionProblemsRawFn func(userSub string) ([]repository.SessionProblemRow, error)
}

func (m *mockMypageRepo) GetSessionProblemsRaw(userSub string) ([]repository.SessionProblemRow, error) {
	return m.getSessionProblemsRawFn(userSub)
}

// --- GetUserData ---

func TestGetUserData_NoSessions(t *testing.T) {
	repo := &mockMypageRepo{
		getSessionProblemsRawFn: func(userSub string) ([]repository.SessionProblemRow, error) {
			return []repository.SessionProblemRow{}, nil
		},
	}
	svc := service.NewMypageService(repo)

	result, err := svc.GetUserData(&model.User{Sub: "sub-1", UserName: "TestUser"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.UserName != "TestUser" {
		t.Errorf("expected UserName = TestUser, got %s", result.UserName)
	}
	if len(result.TestSessDtos) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(result.TestSessDtos))
	}
}

func TestGetUserData_SingleSession_CountsCorrect(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	repo := &mockMypageRepo{
		getSessionProblemsRawFn: func(userSub string) ([]repository.SessionProblemRow, error) {
			return []repository.SessionProblemRow{
				{SessionID: 1, StartTime: now, IsCorrect: true, CategoryName: "足し算"},
				{SessionID: 1, StartTime: now, IsCorrect: false, CategoryName: "引き算"},
				{SessionID: 1, StartTime: now, IsCorrect: true, CategoryName: "掛け算"},
			}, nil
		},
	}
	svc := service.NewMypageService(repo)

	result, err := svc.GetUserData(&model.User{Sub: "sub-1", UserName: "TestUser"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.TestSessDtos) != 1 {
		t.Fatalf("expected 1 session, got %d", len(result.TestSessDtos))
	}
	sess := result.TestSessDtos[0]
	if sess.Total != 3 {
		t.Errorf("expected Total = 3, got %d", sess.Total)
	}
	if sess.CorrectCount != 2 {
		t.Errorf("expected CorrectCount = 2, got %d", sess.CorrectCount)
	}
	if len(sess.CategoryDtos) != 3 {
		t.Errorf("expected 3 category DTOs, got %d", len(sess.CategoryDtos))
	}
	if len(sess.WeakCategories) != 1 || sess.WeakCategories[0] != "引き算" {
		t.Errorf("unexpected WeakCategories: %v", sess.WeakCategories)
	}
}

func TestGetUserData_MultipleSessions_OrderPreserved(t *testing.T) {
	now := time.Now()

	repo := &mockMypageRepo{
		getSessionProblemsRawFn: func(userSub string) ([]repository.SessionProblemRow, error) {
			return []repository.SessionProblemRow{
				{SessionID: 2, StartTime: now, IsCorrect: true, CategoryName: "足し算"},
				{SessionID: 1, StartTime: now, IsCorrect: false, CategoryName: "引き算"},
			}, nil
		},
	}
	svc := service.NewMypageService(repo)

	result, err := svc.GetUserData(&model.User{Sub: "sub-1", UserName: "TestUser"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.TestSessDtos) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(result.TestSessDtos))
	}
	if result.TestSessDtos[0].SessionID != 2 {
		t.Errorf("expected first session ID = 2, got %d", result.TestSessDtos[0].SessionID)
	}
	if result.TestSessDtos[1].SessionID != 1 {
		t.Errorf("expected second session ID = 1, got %d", result.TestSessDtos[1].SessionID)
	}
}

func TestGetUserData_GetSessionProblemsRawError(t *testing.T) {
	repo := &mockMypageRepo{
		getSessionProblemsRawFn: func(userSub string) ([]repository.SessionProblemRow, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewMypageService(repo)

	_, err := svc.GetUserData(&model.User{Sub: "sub-1"})
	if err == nil {
		t.Error("expected error, got nil")
	}
}
