package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
)

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenUseCase interface {
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error)
	CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error)
	ExtractIDFromToken(requestToken string, secret string) (id string, err error)
}
