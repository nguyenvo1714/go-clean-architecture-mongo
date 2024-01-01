package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/utils"
	"time"
)

type RefreshTokenInterceptor struct {
	UserRepository domain.UserRepository
	ContextTimeout time.Duration
}

func (rti *RefreshTokenInterceptor) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	c, cancel := context.WithTimeout(ctx, rti.ContextTimeout)
	defer cancel()

	return rti.UserRepository.GetByID(c, id)
}

func (rti *RefreshTokenInterceptor) CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error) {
	accessToken, err = utils.CreateAccessToken(user, secret, exp)

	return
}

func (rti *RefreshTokenInterceptor) CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error) {
	refreshToken, err = utils.CreateRefreshToken(user, secret, exp)

	return
}

func (rti *RefreshTokenInterceptor) ExtractIDFromToken(requestToken string, secret string) (id string, err error) {
	id, err = utils.ExtractIDFromToken(requestToken, secret)

	return
}
