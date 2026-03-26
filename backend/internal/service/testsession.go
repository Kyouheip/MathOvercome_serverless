package service

import (
	"fmt"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
)

type TestSessionService struct {
	repo repository.TestSessionRepo
}

func NewTestSessionService(r repository.TestSessionRepo) *TestSessionService {
	return &TestSessionService{repo: r}
}

func (s *TestSessionService) CreateTestSess(userSub string, includeIntegers bool) (*model.TestSession, error) {
	session := model.TestSession{
		UserID:          userSub,
		IncludeIntegers: includeIntegers,
	}

	if err := s.repo.SaveTestSession(&session); err != nil {
		return nil, fmt.Errorf("save test session: %w", err)
	}

	maxCategory := 6
	if includeIntegers {
		maxCategory = 7
	}
	var categories []int
	for i := 1; i <= maxCategory; i++ {
		categories = append(categories, i)
	}

	problems, err := s.repo.FindProblemsPerCategory(categories, 2)
	if err != nil {
		return nil, fmt.Errorf("find problems: %w", err)
	}

	var sessProbs []model.SessionProblem
	for _, p := range problems {
		sessProbs = append(sessProbs, model.SessionProblem{
			TestSessionID: session.ID,
			ProblemID:     p.ID,
		})
	}

	if err := s.repo.SaveSessionProblems(sessProbs); err != nil {
		return nil, fmt.Errorf("save session problems: %w", err)
	}

	session.SessionProblems = sessProbs
	return &session, nil
}

func (s *TestSessionService) GetProblem(sessionID uint64, idx int) (*dto.SessionProblem, error) {
	total, err := s.repo.CountSessionProblems(sessionID)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= int(total) {
		return nil, apperr.ErrOutOfRange
	}

	sp, err := s.repo.FindSessionProblemByIdx(sessionID, idx)
	if err != nil {
		return nil, apperr.ErrNotFound
	}

	var choices []dto.Choice
	for _, c := range sp.Problem.Choices {
		choices = append(choices, dto.Choice{
			ID:         int64(c.ID),
			ChoiceText: c.ChoiceText,
		})
	}

	var selectedChoiceID *int64
	if sp.SelectedChoiceID != nil {
		id := int64(*sp.SelectedChoiceID)
		selectedChoiceID = &id
	}

	return &dto.SessionProblem{
		ID:         int64(sp.ID),
		Question:   sp.Problem.Question,
		Choices:    choices,
		Hint:       sp.Problem.Hint,
		SelectedID: selectedChoiceID,
		Total:      int(total),
	}, nil
}

func (s *TestSessionService) SubmitAnswer(sessionID uint64, idx int, choiceID *int64) error {
	if choiceID == nil {
		return nil
	}

	sps, err := s.repo.FindSessionProblemsBySessionID(sessionID)
	if err != nil {
		return err
	}
	if idx < 0 || idx >= len(sps) {
		return apperr.ErrOutOfRange
	}

	sp := sps[idx]
	choice, err := s.repo.FindChoiceByProblemAndChoiceID(sp.ProblemID, uint64(*choiceID))
	if err != nil {
		return apperr.ErrNotFound
	}

	sp.SelectedChoiceID = &choice.ID
	sp.IsCorrect = &choice.IsCorrect
	return s.repo.SaveSessionProblem(&sp)
}
