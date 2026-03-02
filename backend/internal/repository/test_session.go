package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

func (r *Repository) FindTestSessionByID(id uint64) (*model.TestSession, error) {
	var ts model.TestSession
	err := r.db.Preload("User").First(&ts, id).Error
	return &ts, err
}

func (r *Repository) SaveTestSession(session *model.TestSession) error {
	return r.db.Create(session).Error
}
