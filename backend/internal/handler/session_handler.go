package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
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
	user := getUserFromSession(c)
	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	includeIntegers, _ := strconv.ParseBool(c.DefaultQuery("includeIntegers", "false"))

	testSess, err := h.testSessService.CreateTestSess(user, includeIntegers)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(c)
	session.Set("currentSessionId", testSess.ID)
	session.Save()

	c.Status(http.StatusCreated)
}

// GET /session/current/problems/:idx
func (h *SessionHandler) ViewOneProblem(c *gin.Context) {
	user := getUserFromSession(c)
	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	sessionID, ok := getCurrentSessionID(c)
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
	user := getUserFromSession(c)
	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	sessionID, ok := getCurrentSessionID(c)
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
	user := getUserFromSession(c)
	if user == nil {
		c.String(http.StatusUnauthorized, "NOT_LOGIN")
		return
	}

	result, err := h.mypageService.GetUserData(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}

func getUserFromSession(c *gin.Context) *model.User {
	session := sessions.Default(c)
	val := session.Get("user")
	if val == nil {
		return nil
	}
	user, ok := val.(*model.User)
	if !ok {
		return nil
	}
	return user
}

func getCurrentSessionID(c *gin.Context) (uint64, bool) {
	session := sessions.Default(c)
	val := session.Get("currentSessionId")
	if val == nil {
		return 0, false
	}
	id, ok := val.(uint64)
	return id, ok
}
