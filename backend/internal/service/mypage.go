package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
)

type MypageService struct {
	repo repository.MypageRepo
}

func NewMypageService(r repository.MypageRepo) *MypageService {
	return &MypageService{repo: r}
}

func (s *MypageService) GetUserData(user *model.User) (*dto.User, error) {
	rows, err := s.repo.GetSessionProblemsRaw(user.Sub)
	if err != nil {
		return nil, fmt.Errorf("get session problems: %w", err)
	}

	jst := time.FixedZone("JST", 9*60*60)
	sessionMap := make(map[uint64]*dto.TestSession)
	var sessionIDs []uint64

	// catStats[sessionID][categoryName] = {total, correct}
	type catStat struct{ total, correct int }
	catStats := make(map[uint64]map[string]*catStat)
	catOrder := make(map[uint64][]string)

	for _, row := range rows {
		if _, exists := sessionMap[row.SessionID]; !exists {
			d := &dto.TestSession{
				SessionID: int64(row.SessionID),
				StartTime: row.StartTime.In(jst).Format("2006-01-02 15:04:05"),
			}
			sessionMap[row.SessionID] = d
			sessionIDs = append(sessionIDs, row.SessionID)
			catStats[row.SessionID] = make(map[string]*catStat)
		}

		sessionMap[row.SessionID].ProbCategoryDtos = append(
			sessionMap[row.SessionID].ProbCategoryDtos,
			dto.ProblemCategory{
				IsCorrect:    &row.IsCorrect,
				CategoryName: row.CategoryName,
			},
		)

		sm := catStats[row.SessionID]
		if _, exists := sm[row.CategoryName]; !exists {
			sm[row.CategoryName] = &catStat{}
			catOrder[row.SessionID] = append(catOrder[row.SessionID], row.CategoryName)
		}
		sm[row.CategoryName].total++
		if row.IsCorrect {
			sm[row.CategoryName].correct++
		}
	}

	finalSessions := make([]dto.TestSession, 0)
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

		for _, name := range catOrder[id] {
			st := catStats[id][name]
			sess.CategoryDtos = append(sess.CategoryDtos, dto.Category{
				CategoryName: name,
				Total:        st.total,
				CorrectCount: st.correct,
			})
		}

		type weakCat struct {
			name string
			rate float64
		}
		var weaks []weakCat
		for _, name := range catOrder[id] {
			st := catStats[id][name]
			if st.total == 0 {
				continue
			}
			rate := float64(st.correct) / float64(st.total)
			if rate < 0.5 {
				weaks = append(weaks, weakCat{name: name, rate: rate})
			}
		}
		sort.Slice(weaks, func(i, j int) bool {
			return weaks[i].rate < weaks[j].rate
		})
		for i := 0; i < len(weaks) && i < 2; i++ {
			sess.WeakCategories = append(sess.WeakCategories, weaks[i].name)
		}

		finalSessions = append(finalSessions, *sess)
	}

	sort.Slice(finalSessions, func(i, j int) bool {
		return finalSessions[i].StartTime > finalSessions[j].StartTime
	})

	return &dto.User{
		UserName:     user.UserName,
		TestSessDtos: finalSessions,
	}, nil
}
