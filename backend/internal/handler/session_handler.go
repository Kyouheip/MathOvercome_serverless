package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

type SessionHandler struct {
	testSessService service.TestSessionServicer
	mypageService   service.MypageServicer
}

func NewSessionHandler(ts service.TestSessionServicer, ms service.MypageServicer) *SessionHandler {
	return &SessionHandler{testSessService: ts, mypageService: ms}
}

// POST /session/test
func (h *SessionHandler) CreateTestSess(c *gin.Context) {
	userSub := c.GetHeader("X-User-Sub")
	if userSub == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	includeIntegers, _ := strconv.ParseBool(c.DefaultQuery("includeIntegers", "false"))

	testSess, err := h.testSessService.CreateTestSess(userSub, includeIntegers)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"sessionId": testSess.ID})
}

// GET /session/current/problems/:idx
func (h *SessionHandler) ViewOneProblem(c *gin.Context) {
	if c.GetHeader("X-User-Sub") == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	sessionID, ok := getSessionIDFromQuery(c)
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	idx, _ := strconv.Atoi(c.Param("idx"))

	problem, err := h.testSessService.GetProblem(sessionID, idx)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrOutOfRange):
			c.Status(http.StatusBadRequest)
		case errors.Is(err, apperr.ErrNotFound):
			c.Status(http.StatusNotFound)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, problem)
}

// POST /session/current/problems/:idx/answer
func (h *SessionHandler) SubmitAnswer(c *gin.Context) {
	if c.GetHeader("X-User-Sub") == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	sessionID, ok := getSessionIDFromQuery(c)
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	idx, _ := strconv.Atoi(c.Param("idx"))

	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.testSessService.SubmitAnswer(sessionID, idx, req.SelectedChoiceID); err != nil {
		switch {
		case errors.Is(err, apperr.ErrOutOfRange), errors.Is(err, apperr.ErrNotFound):
			c.Status(http.StatusBadRequest)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /session/mypage
func (h *SessionHandler) GetMypage(c *gin.Context) {
	userSub := c.GetHeader("X-User-Sub")
	if userSub == "" {
		c.String(http.StatusUnauthorized, "NOT_LOGIN")
		return
	}

	user := &model.User{
		Sub:      userSub,
		UserName: c.GetHeader("X-User-Name"),
	}

	result, err := h.mypageService.GetUserData(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}

func getSessionIDFromQuery(c *gin.Context) (uint64, bool) {
	s := c.Query("sessionId")
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}
