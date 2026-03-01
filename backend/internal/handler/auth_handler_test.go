package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// mockLoginService は service.LoginServicer のテスト用実装。
type mockLoginService struct {
	authenticateFn     func(req dto.LoginRequest) (*model.User, error)
	validateRegisterFn func(req dto.RegisterRequest) error
	createUserFn       func(req dto.RegisterRequest) error
}

func (m *mockLoginService) Authenticate(req dto.LoginRequest) (*model.User, error) {
	return m.authenticateFn(req)
}

func (m *mockLoginService) ValidateRegister(req dto.RegisterRequest) error {
	return m.validateRegisterFn(req)
}

func (m *mockLoginService) CreateUser(req dto.RegisterRequest) error {
	return m.createUserFn(req)
}

// newAuthEngine はテスト用の gin エンジンを作成する。
func newAuthEngine(svc service.LoginServicer) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("session", store))
	h := handler.NewAuthHandler(svc)
	r.GET("/auth/ping", h.Ping)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
	r.POST("/auth/register", h.Register)
	return r
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return bytes.NewReader(b)
}

// --- Ping ---

func TestPing(t *testing.T) {
	r := newAuthEngine(&mockLoginService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/ping", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- Login ---

func TestLogin_Success(t *testing.T) {
	svc := &mockLoginService{
		authenticateFn: func(req dto.LoginRequest) (*model.User, error) {
			return &model.User{ID: 1, UserName: "Test", UserID: req.UserID}, nil
		},
	}
	r := newAuthEngine(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, dto.LoginRequest{UserID: "testuser", Password: "pass123"}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	r := newAuthEngine(&mockLoginService{})
	// password フィールドなし
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, map[string]string{"userId": "testuser"}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	svc := &mockLoginService{
		authenticateFn: func(req dto.LoginRequest) (*model.User, error) {
			return nil, fmt.Errorf("authenticate: %w", apperr.ErrInvalidCredentials)
		},
	}
	r := newAuthEngine(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", jsonBody(t, dto.LoginRequest{UserID: "baduser", Password: "wrong"}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- Logout ---

func TestLogout(t *testing.T) {
	r := newAuthEngine(&mockLoginService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

// --- Register ---

func TestRegister_Success(t *testing.T) {
	svc := &mockLoginService{
		validateRegisterFn: func(req dto.RegisterRequest) error { return nil },
		createUserFn:       func(req dto.RegisterRequest) error { return nil },
	}
	r := newAuthEngine(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "testuser1",
		Password1: "pass123",
		Password2: "pass123",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestRegister_UserIDTooShort(t *testing.T) {
	r := newAuthEngine(&mockLoginService{})
	// UserID が 6 文字未満 → バリデーションエラー
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "abc",
		Password1: "pass123",
		Password2: "pass123",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	svc := &mockLoginService{
		validateRegisterFn: func(req dto.RegisterRequest) error {
			return fmt.Errorf("validate: %w", apperr.ErrPasswordMismatch)
		},
	}
	r := newAuthEngine(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "testuser1",
		Password1: "pass123",
		Password2: "pass124",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	svc := &mockLoginService{
		validateRegisterFn: func(req dto.RegisterRequest) error {
			return fmt.Errorf("validate: %w", apperr.ErrUserAlreadyExists)
		},
	}
	r := newAuthEngine(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", jsonBody(t, dto.RegisterRequest{
		UserName:  "Test User",
		UserID:    "existing1",
		Password1: "pass123",
		Password2: "pass123",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
