package service

import (
	"fmt"

	"github.com/Kyouheip/MathOvercome_serverless/internal/apperr"
	"github.com/Kyouheip/MathOvercome_serverless/internal/dto"
	"github.com/Kyouheip/MathOvercome_serverless/internal/model"
	"github.com/Kyouheip/MathOvercome_serverless/internal/repository"
)

type LoginService struct {
	repo *repository.Repository
}

func NewLoginService(r *repository.Repository) *LoginService {
	return &LoginService{repo: r}
}

func (s *LoginService) Authenticate(req dto.LoginRequest) (*model.User, error) {
	user, err := s.repo.FindUserByUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", apperr.ErrInvalidCredentials)
	}

	if user.Password != req.Password {
		return nil, fmt.Errorf("authenticate: %w", apperr.ErrInvalidCredentials)
	}

	return user, nil
}

func (s *LoginService) ValidateRegister(req dto.RegisterRequest) error {
	if req.Password1 != req.Password2 {
		return fmt.Errorf("validate register: %w", apperr.ErrPasswordMismatch)
	}

	_, err := s.repo.FindUserByUserID(req.UserID)
	if err == nil {
		return fmt.Errorf("validate register: %w", apperr.ErrUserAlreadyExists)
	}

	return nil
}

func (s *LoginService) CreateUser(req dto.RegisterRequest) error {
	user := model.User{
		UserName: req.UserName,
		UserID:   req.UserID,
		Password: req.Password1,
	}
	if err := s.repo.SaveUser(&user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}
