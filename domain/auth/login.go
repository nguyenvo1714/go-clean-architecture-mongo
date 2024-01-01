package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUseCase interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error)
	CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error)
}
