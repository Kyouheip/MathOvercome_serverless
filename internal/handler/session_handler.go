package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
	"github.com/Kyouheip/MathOvercome_serverless/internal/service"
)

type SessionHandler struct {
	testSessService *service.TestSessionService
	mypageService   *service.MypageService
	repo            *repository.Repository
}

func NewSessionHandler(ts *service.TestSessionService, ms *service.MypageService, r *repository.Repository) *SessionHandler {
	return &SessionHandler{testSessService: ts, mypageService: ms, repo: r}
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

	ts, err := h.repo.FindTestSessionByID(sessionID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if ts.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	total, err := h.repo.CountSessionProblems(sessionID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if idx < 0 || idx >= int(total) {
		c.Status(http.StatusBadRequest)
		return
	}

	sp, err := h.repo.FindSessionProblemByIdx(sessionID, idx)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var choices []dto.Choice
	for _, choice := range sp.Problem.Choices {
		choices = append(choices, dto.Choice{
			ID:         int64(choice.ID),
			ChoiceText: choice.ChoiceText,
		})
	}

	var selectedChoiceID *int64
	if sp.SelectedChoiceID != nil {
		id := int64(*sp.SelectedChoiceID)
		selectedChoiceID = &id
	}

	c.JSON(http.StatusOK, dto.SessionProblem{
		ID:         int64(sp.ID),
		Question:   sp.Problem.Question,
		Choices:    choices,
		Hint:       sp.Problem.Hint,
		SelectedID: selectedChoiceID,
		Total:      int(total),
	})
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

	ts, err := h.repo.FindTestSessionByID(sessionID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if ts.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	sps, err := h.repo.FindSessionProblemsBySessionID(sessionID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if idx < 0 || idx >= len(sps) {
		c.Status(http.StatusBadRequest)
		return
	}

	var req dto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// 未解答(null)の場合は処理なし
	if req.SelectedChoiceID != nil {
		sp := sps[idx]
		selectedChoice, err := h.repo.FindChoiceByID(uint64(*req.SelectedChoiceID))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		sp.SelectedChoiceID = &selectedChoice.ID
		sp.IsCorrect = &selectedChoice.IsCorrect
		if err := h.repo.SaveSessionProblem(&sp); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
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
