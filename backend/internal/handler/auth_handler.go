package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

type AuthHandler struct {
	loginService service.LoginServicer
}

func NewAuthHandler(s service.LoginServicer) *AuthHandler {
	return &AuthHandler{loginService: s}
}

// GET /auth/ping
func (h *AuthHandler) Ping(c *gin.Context) {
	c.Status(http.StatusOK)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := collectValidationErrors(err)
		if len(msgs) == 0 {
			msgs = append(msgs, "IDまたはパスワードを入力してください")
		}
		c.String(http.StatusBadRequest, strings.Join(msgs, "\n"))
		return
	}

	user, err := h.loginService.Authenticate(req)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			c.String(http.StatusBadRequest, "IDまたはパスワードが間違っています")
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(c)
	session.Set("user", user)
	session.Save()

	c.Status(http.StatusOK)
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Status(http.StatusNoContent)
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msgs := collectRegisterErrors(err)
		c.String(http.StatusBadRequest, strings.Join(msgs, "\n"))
		return
	}

	if err := h.loginService.ValidateRegister(req); err != nil {
		switch {
		case errors.Is(err, apperr.ErrPasswordMismatch):
			c.String(http.StatusBadRequest, "パスワードが一致しません")
		case errors.Is(err, apperr.ErrUserAlreadyExists):
			c.String(http.StatusBadRequest, "すでに利用されているIDです")
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	if err := h.loginService.CreateUser(req); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

func collectValidationErrors(err error) []string {
	var msgs []string
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			msgs = append(msgs, fe.Error())
		}
	}
	return msgs
}

func collectRegisterErrors(err error) []string {
	fieldNameMap := map[string]string{
		"Password2": "パスワード確認",
		"Password1": "パスワード",
		"UserName":  "名前",
		"UserID":    "ID",
	}

	var msgs []string
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := fe.Field()
			japanField, ok := fieldNameMap[field]
			if !ok {
				japanField = field
			}
			msgs = append(msgs, fmt.Sprintf("%s:%s", japanField, fe.Tag()))
		}
	}
	return msgs
}
