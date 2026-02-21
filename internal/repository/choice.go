package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

func (r *Repository) FindChoicesByProblemID(problemID uint64) ([]model.Choice, error) {
	var choices []model.Choice
	err := r.db.Where("problem_id = ?", problemID).Find(&choices).Error
	return choices, err
}

func (r *Repository) FindChoiceByID(id uint64) (*model.Choice, error) {
	var choice model.Choice
	err := r.db.First(&choice, id).Error
	return &choice, err
}
