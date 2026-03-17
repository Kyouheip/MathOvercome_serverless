package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

// --- モック実装 ---

type mockTestSessionService struct {
	createTestSessFn func(user *model.User, includeIntegers bool) (*model.TestSession, error)
	getProblemFn     func(sessionID uint64, idx int) (*dto.SessionProblem, error)
	submitAnswerFn   func(sessionID uint64, idx int, choiceID *int64) error
}

func (m *mockTestSessionService) CreateTestSess(user *model.User, includeIntegers bool) (*model.TestSession, error) {
	return m.createTestSessFn(user, includeIntegers)
}

func (m *mockTestSessionService) GetProblem(sessionID uint64, idx int) (*dto.SessionProblem, error) {
	return m.getProblemFn(sessionID, idx)
}

func (m *mockTestSessionService) SubmitAnswer(sessionID uint64, idx int, choiceID *int64) error {
	return m.submitAnswerFn(sessionID, idx, choiceID)
}

type mockMypageService struct {
	getUserDataFn func(user *model.User) (*dto.User, error)
}

func (m *mockMypageService) GetUserData(user *model.User) (*dto.User, error) {
	return m.getUserDataFn(user)
}

// newSessionEngine はテスト用エンジンを作成する。
func newSessionEngine(
	ts service.TestSessionServicer,
	ms service.MypageServicer,
	sessionVals map[string]any,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("session", store))

	if len(sessionVals) > 0 {
		r.Use(func(c *gin.Context) {
			sess := sessions.Default(c)
			for k, v := range sessionVals {
				sess.Set(k, v)
			}
			c.Next()
		})
	}

	h := handler.NewSessionHandler(ts, ms)
	r.POST("/session/test", h.CreateTestSess)
	r.GET("/session/current/problems/:idx", h.ViewOneProblem)
	r.POST("/session/current/problems/:idx/answer", h.SubmitAnswer)
	r.GET("/session/mypage", h.GetMypage)
	return r
}

// --- CreateTestSess ---

func TestCreateTestSess_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateTestSess_Success(t *testing.T) {
	user := &model.User{ID: 1}
	ts := &mockTestSessionService{
		createTestSessFn: func(u *model.User, includeIntegers bool) (*model.TestSession, error) {
			return &model.TestSession{ID: 42, UserID: u.ID}, nil
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCreateTestSess_ServiceError(t *testing.T) {
	user := &model.User{ID: 1}
	ts := &mockTestSessionService{
		createTestSessFn: func(u *model.User, includeIntegers bool) (*model.TestSession, error) {
			return nil, errors.New("db error")
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- ViewOneProblem ---

func TestViewOneProblem_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestViewOneProblem_NoCurrentSession(t *testing.T) {
	user := &model.User{ID: 1}
	r := newSessionEngine(nil, nil, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_OutOfRange(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	ts := &mockTestSessionService{
		getProblemFn: func(sID uint64, idx int) (*dto.SessionProblem, error) {
			return nil, apperr.ErrOutOfRange
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_Success(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	choiceID := int64(5)

	ts := &mockTestSessionService{
		getProblemFn: func(sID uint64, idx int) (*dto.SessionProblem, error) {
			return &dto.SessionProblem{
				ID:       1,
				Question: "1+1=?",
				Hint:     "2です",
				Choices: []dto.Choice{
					{ID: 5, ChoiceText: "2"},
					{ID: 6, ChoiceText: "3"},
				},
				SelectedID: &choiceID,
				Total:      5,
			}, nil
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp dto.SessionProblem
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Question != "1+1=?" {
		t.Errorf("expected Question = 1+1=?, got %s", resp.Question)
	}
	if len(resp.Choices) != 2 {
		t.Errorf("expected 2 choices, got %d", len(resp.Choices))
	}
	if resp.Total != 5 {
		t.Errorf("expected Total = 5, got %d", resp.Total)
	}
	if resp.SelectedID == nil || *resp.SelectedID != 5 {
		t.Errorf("expected SelectedID = 5, got %v", resp.SelectedID)
	}
}

// --- SubmitAnswer ---

func TestSubmitAnswer_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestSubmitAnswer_NullAnswer(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return nil
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: nil})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestSubmitAnswer_WithChoice(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return nil
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	choiceID := int64(5)
	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: &choiceID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestSubmitAnswer_OutOfRange(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return apperr.ErrOutOfRange
		},
	}
	r := newSessionEngine(ts, nil, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: nil})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/5/answer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GetMypage ---

func TestGetMypage_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetMypage_Success(t *testing.T) {
	user := &model.User{ID: 1, UserName: "TestUser"}
	ms := &mockMypageService{
		getUserDataFn: func(u *model.User) (*dto.User, error) {
			return &dto.User{UserName: u.UserName, TestSessDtos: []dto.TestSession{}}, nil
		},
	}
	r := newSessionEngine(nil, ms, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp dto.User
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.UserName != "TestUser" {
		t.Errorf("expected UserName = TestUser, got %s", resp.UserName)
	}
}

func TestGetMypage_ServiceError(t *testing.T) {
	user := &model.User{ID: 1}
	ms := &mockMypageService{
		getUserDataFn: func(u *model.User) (*dto.User, error) {
			return nil, errors.New("db error")
		},
	}
	r := newSessionEngine(nil, ms, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
