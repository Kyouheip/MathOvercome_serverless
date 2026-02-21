package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

// Java の sp.getProblem().getChoices() に相当するため Problem.Choices も Preload
func (r *Repository) FindSessionProblemByIdx(sessionID uint64, idx int) (*model.SessionProblem, error) {
	var sp model.SessionProblem
	err := r.db.
		Preload("Problem").
		Preload("Problem.Choices").
		Where("session_id = ?", sessionID).
		Order("id").
		Offset(idx).
		Limit(1).
		First(&sp).Error
	return &sp, err
}

func (r *Repository) CountSessionProblems(sessionID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.SessionProblem{}).Where("session_id = ?", sessionID).Count(&count).Error
	return count, err
}

func (r *Repository) FindSessionProblemsBySessionID(sessionID uint64) ([]model.SessionProblem, error) {
	var sps []model.SessionProblem
	err := r.db.Where("session_id = ?", sessionID).Order("id").Find(&sps).Error
	return sps, err
}

func (r *Repository) SaveSessionProblem(sp *model.SessionProblem) error {
	return r.db.Save(sp).Error
}

func (r *Repository) SaveSessionProblems(sps []model.SessionProblem) error {
	return r.db.Create(&sps).Error
}
