package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

// --- モック実装 ---

type mockTestSessionService struct {
	createTestSessFn func(userSub string, includeIntegers bool) (*model.TestSession, error)
	getProblemFn     func(sessionID uint64, idx int) (*dto.SessionProblem, error)
	submitAnswerFn   func(sessionID uint64, idx int, choiceID *int64) error
}

func (m *mockTestSessionService) CreateTestSess(userSub string, includeIntegers bool) (*model.TestSession, error) {
	return m.createTestSessFn(userSub, includeIntegers)
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
// userSub が空でない場合、X-User-Sub ヘッダーをリクエストにセットするミドルウェアを追加する。
func newSessionEngine(
	ts service.TestSessionServicer,
	ms service.MypageServicer,
	userSub string,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	h := handler.NewSessionHandler(ts, ms)
	r.POST("/session/test", h.CreateTestSess)
	r.GET("/session/current/problems/:idx", h.ViewOneProblem)
	r.POST("/session/current/problems/:idx/answer", h.SubmitAnswer)
	r.GET("/session/mypage", h.GetMypage)
	return r
}

func addUserSub(req *http.Request, sub string) {
	req.Header.Set("X-User-Sub", sub)
}

// --- CreateTestSess ---

func TestCreateTestSess_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateTestSess_Success(t *testing.T) {
	ts := &mockTestSessionService{
		createTestSessFn: func(userSub string, includeIntegers bool) (*model.TestSession, error) {
			return &model.TestSession{ID: 42, UserID: userSub}, nil
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCreateTestSess_ServiceError(t *testing.T) {
	ts := &mockTestSessionService{
		createTestSessFn: func(userSub string, includeIntegers bool) (*model.TestSession, error) {
			return nil, errors.New("db error")
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- ViewOneProblem ---

func TestViewOneProblem_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestViewOneProblem_NoSessionID(t *testing.T) {
	r := newSessionEngine(nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_OutOfRange(t *testing.T) {
	ts := &mockTestSessionService{
		getProblemFn: func(sID uint64, idx int) (*dto.SessionProblem, error) {
			return nil, apperr.ErrOutOfRange
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/10?sessionId=10", nil)
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_Success(t *testing.T) {
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
	r := newSessionEngine(ts, nil, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0?sessionId=10", nil)
	addUserSub(req, "sub-1")
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
	r := newSessionEngine(nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestSubmitAnswer_NullAnswer(t *testing.T) {
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return nil
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: nil})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer?sessionId=10", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestSubmitAnswer_WithChoice(t *testing.T) {
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return nil
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	choiceID := int64(5)
	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: &choiceID})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/0/answer?sessionId=10", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestSubmitAnswer_OutOfRange(t *testing.T) {
	ts := &mockTestSessionService{
		submitAnswerFn: func(sID uint64, idx int, choiceID *int64) error {
			return apperr.ErrOutOfRange
		},
	}
	r := newSessionEngine(ts, nil, "sub-1")

	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: nil})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/5/answer?sessionId=10", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GetMypage ---

func TestGetMypage_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, "")
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetMypage_Success(t *testing.T) {
	ms := &mockMypageService{
		getUserDataFn: func(u *model.User) (*dto.User, error) {
			return &dto.User{UserName: u.UserName, TestSessDtos: []dto.TestSession{}}, nil
		},
	}
	r := newSessionEngine(nil, ms, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	req.Header.Set("X-User-Sub", "sub-1")
	req.Header.Set("X-User-Name", "TestUser")
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
	ms := &mockMypageService{
		getUserDataFn: func(u *model.User) (*dto.User, error) {
			return nil, errors.New("db error")
		},
	}
	r := newSessionEngine(nil, ms, "sub-1")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	addUserSub(req, "sub-1")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
