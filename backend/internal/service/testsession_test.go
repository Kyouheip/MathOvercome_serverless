package service_test

import (
	"errors"
	"testing"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

// mockTestSessionRepo は repository.TestSessionRepo のテスト用実装。
type mockTestSessionRepo struct {
	saveTestSessionFn         func(session *model.TestSession) error
	findProblemsPerCategoryFn func(categoryIDs []int, countPerCategory int) ([]model.Problem, error)
	saveSessionProblemsFn     func(sps []model.SessionProblem) error
}

func (m *mockTestSessionRepo) SaveTestSession(session *model.TestSession) error {
	return m.saveTestSessionFn(session)
}

func (m *mockTestSessionRepo) FindProblemsPerCategory(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
	return m.findProblemsPerCategoryFn(categoryIDs, countPerCategory)
}

func (m *mockTestSessionRepo) SaveSessionProblems(sps []model.SessionProblem) error {
	return m.saveSessionProblemsFn(sps)
}

func makeProblems(n int) []model.Problem {
	probs := make([]model.Problem, n)
	for i := range probs {
		probs[i] = model.Problem{ID: uint64(i + 1)}
	}
	return probs
}

// --- CreateTestSess ---

func TestCreateTestSess_Success(t *testing.T) {
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			session.ID = 1
			return nil
		},
		findProblemsPerCategoryFn: func(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
			return makeProblems(len(categoryIDs) * countPerCategory), nil
		},
		saveSessionProblemsFn: func(sps []model.SessionProblem) error { return nil },
	}
	svc := service.NewTestSessionService(repo)

	sess, err := svc.CreateTestSess(&model.User{ID: 1}, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 6カテゴリ × 2問 = 12問
	if len(sess.SessionProblems) != 12 {
		t.Errorf("expected 12 session problems, got %d", len(sess.SessionProblems))
	}
	if sess.IncludeIntegers != false {
		t.Error("expected IncludeIntegers = false")
	}
}

func TestCreateTestSess_WithIntegers(t *testing.T) {
	var gotCategoryIDs []int
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			session.ID = 1
			return nil
		},
		findProblemsPerCategoryFn: func(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
			gotCategoryIDs = categoryIDs
			return makeProblems(len(categoryIDs) * countPerCategory), nil
		},
		saveSessionProblemsFn: func(sps []model.SessionProblem) error { return nil },
	}
	svc := service.NewTestSessionService(repo)

	sess, err := svc.CreateTestSess(&model.User{ID: 1}, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// 7カテゴリ × 2問 = 14問
	if len(sess.SessionProblems) != 14 {
		t.Errorf("expected 14 session problems, got %d", len(sess.SessionProblems))
	}
	if len(gotCategoryIDs) != 7 {
		t.Errorf("expected 7 categories, got %d", len(gotCategoryIDs))
	}
}

func TestCreateTestSess_SessionProblemsLinkedToSession(t *testing.T) {
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			session.ID = 99
			return nil
		},
		findProblemsPerCategoryFn: func(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
			return makeProblems(len(categoryIDs) * countPerCategory), nil
		},
		saveSessionProblemsFn: func(sps []model.SessionProblem) error { return nil },
	}
	svc := service.NewTestSessionService(repo)

	sess, err := svc.CreateTestSess(&model.User{ID: 5}, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for _, sp := range sess.SessionProblems {
		if sp.TestSessionID != 99 {
			t.Errorf("expected TestSessionID = 99, got %d", sp.TestSessionID)
		}
	}
}

func TestCreateTestSess_SaveSessionError(t *testing.T) {
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			return errors.New("db error")
		},
	}
	svc := service.NewTestSessionService(repo)

	_, err := svc.CreateTestSess(&model.User{ID: 1}, false)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateTestSess_FindProblemsError(t *testing.T) {
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			session.ID = 1
			return nil
		},
		findProblemsPerCategoryFn: func(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTestSessionService(repo)

	_, err := svc.CreateTestSess(&model.User{ID: 1}, false)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCreateTestSess_SaveSessionProblemsError(t *testing.T) {
	repo := &mockTestSessionRepo{
		saveTestSessionFn: func(session *model.TestSession) error {
			session.ID = 1
			return nil
		},
		findProblemsPerCategoryFn: func(categoryIDs []int, countPerCategory int) ([]model.Problem, error) {
			return makeProblems(len(categoryIDs) * countPerCategory), nil
		},
		saveSessionProblemsFn: func(sps []model.SessionProblem) error {
			return errors.New("db error")
		},
	}
	svc := service.NewTestSessionService(repo)

	_, err := svc.CreateTestSess(&model.User{ID: 1}, false)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
