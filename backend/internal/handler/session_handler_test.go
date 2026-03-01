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

	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/handler"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

// --- モック実装 ---

type mockTestSessionService struct {
	createTestSessFn func(user *model.User, includeIntegers bool) (*model.TestSession, error)
}

func (m *mockTestSessionService) CreateTestSess(user *model.User, includeIntegers bool) (*model.TestSession, error) {
	return m.createTestSessFn(user, includeIntegers)
}

type mockMypageService struct {
	getUserDataFn func(user *model.User) (*dto.User, error)
}

func (m *mockMypageService) GetUserData(user *model.User) (*dto.User, error) {
	return m.getUserDataFn(user)
}

type mockSessionRepo struct {
	findTestSessionByIDFn           func(id uint64) (*model.TestSession, error)
	countSessionProblemsFn          func(sessionID uint64) (int64, error)
	findSessionProblemByIdxFn       func(sessionID uint64, idx int) (*model.SessionProblem, error)
	findSessionProblemsBySessionIDFn func(sessionID uint64) ([]model.SessionProblem, error)
	findChoiceByIDFn                func(id uint64) (*model.Choice, error)
	saveSessionProblemFn            func(sp *model.SessionProblem) error
}

func (m *mockSessionRepo) FindTestSessionByID(id uint64) (*model.TestSession, error) {
	return m.findTestSessionByIDFn(id)
}

func (m *mockSessionRepo) CountSessionProblems(sessionID uint64) (int64, error) {
	return m.countSessionProblemsFn(sessionID)
}

func (m *mockSessionRepo) FindSessionProblemByIdx(sessionID uint64, idx int) (*model.SessionProblem, error) {
	return m.findSessionProblemByIdxFn(sessionID, idx)
}

func (m *mockSessionRepo) FindSessionProblemsBySessionID(sessionID uint64) ([]model.SessionProblem, error) {
	return m.findSessionProblemsBySessionIDFn(sessionID)
}

func (m *mockSessionRepo) FindChoiceByID(id uint64) (*model.Choice, error) {
	return m.findChoiceByIDFn(id)
}

func (m *mockSessionRepo) SaveSessionProblem(sp *model.SessionProblem) error {
	return m.saveSessionProblemFn(sp)
}

// newSessionEngine はテスト用エンジンを作成する。
// sessionVals に指定した値はリクエスト処理前にセッションへセットされる。
func newSessionEngine(
	ts service.TestSessionServicer,
	ms service.MypageServicer,
	repo repository.SessionRepo,
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

	h := handler.NewSessionHandler(ts, ms, repo)
	r.POST("/session/test", h.CreateTestSess)
	r.GET("/session/current/problems/:idx", h.ViewOneProblem)
	r.POST("/session/current/problems/:idx/answer", h.SubmitAnswer)
	r.GET("/session/mypage", h.GetMypage)
	return r
}

// --- CreateTestSess ---

func TestCreateTestSess_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, &mockSessionRepo{}, nil)
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
	r := newSessionEngine(ts, nil, &mockSessionRepo{}, map[string]any{"user": user})

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
	r := newSessionEngine(ts, nil, &mockSessionRepo{}, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/session/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- ViewOneProblem ---

func TestViewOneProblem_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, &mockSessionRepo{}, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestViewOneProblem_NoCurrentSession(t *testing.T) {
	user := &model.User{ID: 1}
	// currentSessionId をセットしない
	r := newSessionEngine(nil, nil, &mockSessionRepo{}, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_Forbidden(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			// セッションのオーナーは UserID=99
			return &model.TestSession{ID: id, UserID: 99}, nil
		},
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestViewOneProblem_OutOfRange(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			return &model.TestSession{ID: id, UserID: 1}, nil
		},
		countSessionProblemsFn: func(sID uint64) (int64, error) {
			return 5, nil
		},
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	w := httptest.NewRecorder()
	// idx=10 は範囲外 (total=5)
	req := httptest.NewRequest(http.MethodGet, "/session/current/problems/10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestViewOneProblem_Success(t *testing.T) {
	user := &model.User{ID: 1}
	var sessionID uint64 = 10
	choiceID := uint64(5)
	isCorrect := true

	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			return &model.TestSession{ID: id, UserID: 1}, nil
		},
		countSessionProblemsFn: func(sID uint64) (int64, error) {
			return 5, nil
		},
		findSessionProblemByIdxFn: func(sID uint64, idx int) (*model.SessionProblem, error) {
			return &model.SessionProblem{
				ID: 1,
				Problem: model.Problem{
					Question: "1+1=?",
					Hint:     "2です",
					Choices: []model.Choice{
						{ID: 5, ChoiceText: "2", IsCorrect: true},
						{ID: 6, ChoiceText: "3", IsCorrect: false},
					},
				},
				SelectedChoiceID: &choiceID,
				IsCorrect:        &isCorrect,
			}, nil
		},
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
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
	r := newSessionEngine(nil, nil, &mockSessionRepo{}, nil)
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
	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			return &model.TestSession{ID: id, UserID: 1}, nil
		},
		findSessionProblemsBySessionIDFn: func(sID uint64) ([]model.SessionProblem, error) {
			return []model.SessionProblem{{ID: 1, ProblemID: 1}}, nil
		},
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
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

	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			return &model.TestSession{ID: id, UserID: 1}, nil
		},
		findSessionProblemsBySessionIDFn: func(sID uint64) ([]model.SessionProblem, error) {
			return []model.SessionProblem{{ID: 1, ProblemID: 1}}, nil
		},
		findChoiceByIDFn: func(id uint64) (*model.Choice, error) {
			return &model.Choice{ID: id, IsCorrect: true}, nil
		},
		saveSessionProblemFn: func(sp *model.SessionProblem) error { return nil },
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
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
	repo := &mockSessionRepo{
		findTestSessionByIDFn: func(id uint64) (*model.TestSession, error) {
			return &model.TestSession{ID: id, UserID: 1}, nil
		},
		findSessionProblemsBySessionIDFn: func(sID uint64) ([]model.SessionProblem, error) {
			return []model.SessionProblem{{ID: 1}}, nil // 問題は1問のみ
		},
	}
	r := newSessionEngine(nil, nil, repo, map[string]any{
		"user":             user,
		"currentSessionId": sessionID,
	})

	body, _ := json.Marshal(dto.AnswerRequest{SelectedChoiceID: nil})
	w := httptest.NewRecorder()
	// idx=5 は範囲外
	req := httptest.NewRequest(http.MethodPost, "/session/current/problems/5/answer", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GetMypage ---

func TestGetMypage_Unauthorized(t *testing.T) {
	r := newSessionEngine(nil, nil, &mockSessionRepo{}, nil)
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
	r := newSessionEngine(nil, ms, &mockSessionRepo{}, map[string]any{"user": user})

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
	r := newSessionEngine(nil, ms, &mockSessionRepo{}, map[string]any{"user": user})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/session/mypage", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
