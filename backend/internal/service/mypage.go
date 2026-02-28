package service

import (
	"fmt"
	"time"

	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
)

type MypageService struct {
	repo *repository.Repository
}

func NewMypageService(r *repository.Repository) *MypageService {
	return &MypageService{repo: r}
}

func (s *MypageService) GetUserData(user *model.User) (*dto.User, error) {
	rows, err := s.repo.GetSessionProblemsRaw(user.ID)
	if err != nil {
		return nil, fmt.Errorf("get session problems: %w", err)
	}

	jst := time.FixedZone("JST", 9*60*60)
	sessionMap := make(map[uint64]*dto.TestSession)
	var sessionIDs []uint64

	for _, row := range rows {
		if _, exists := sessionMap[row.SessionID]; !exists {
			d := &dto.TestSession{
				SessionID: int64(row.SessionID),
				StartTime: row.StartTime.In(jst).Format("2006-01-02 15:04:05"),
			}
			sessionMap[row.SessionID] = d
			sessionIDs = append(sessionIDs, row.SessionID)
		}

		sessionMap[row.SessionID].ProbCategoryDtos = append(
			sessionMap[row.SessionID].ProbCategoryDtos,
			dto.ProblemCategory{
				IsCorrect:    &row.IsCorrect,
				CategoryName: row.CategoryName,
			},
		)
	}

	var finalSessions []dto.TestSession
	for _, id := range sessionIDs {
		sess := sessionMap[id]

		total := len(sess.ProbCategoryDtos)
		correct := 0
		for _, p := range sess.ProbCategoryDtos {
			if *p.IsCorrect {
				correct++
			}
		}
		sess.Total = total
		sess.CorrectCount = correct

		stats, err := s.repo.GetCategoryStats(id)
		if err != nil {
			return nil, fmt.Errorf("get category stats for session %d: %w", id, err)
		}
		for _, st := range stats {
			sess.CategoryDtos = append(sess.CategoryDtos, dto.Category{
				CategoryName: st.Name,
				Total:        st.TotalCount,
				CorrectCount: st.CorrectCount,
			})
		}

		weak, err := s.repo.GetWeakCategories(id)
		if err != nil {
			return nil, fmt.Errorf("get weak categories for session %d: %w", id, err)
		}
		sess.WeakCategories = weak

		finalSessions = append(finalSessions, *sess)
	}

	return &dto.User{
		UserName:     user.UserName,
		TestSessDtos: finalSessions,
	}, nil
}
