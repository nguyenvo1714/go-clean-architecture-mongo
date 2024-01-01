package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
)

type Signup struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupUseCase interface {
	Create(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error)
	CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error)
}
