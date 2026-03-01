package service

import (
	"fmt"

	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
)

type TestSessionService struct {
	repo repository.TestSessionRepo
}

func NewTestSessionService(r repository.TestSessionRepo) *TestSessionService {
	return &TestSessionService{repo: r}
}

func (s *TestSessionService) CreateTestSess(user *model.User, includeIntegers bool) (*model.TestSession, error) {
	session := model.TestSession{
		UserID:          user.ID,
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
