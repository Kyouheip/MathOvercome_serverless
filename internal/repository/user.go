package repository

import "github.com/Kyouheip/MathOvercome_serverless/internal/model"

func (r *Repository) FindUserByUserID(userID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (r *Repository) SaveUser(user *model.User) error {
	return r.db.Create(user).Error
}
