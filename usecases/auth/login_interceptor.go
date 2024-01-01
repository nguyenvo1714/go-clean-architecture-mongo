package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/utils"
	"time"
)

type LoginInterceptor struct {
	UserRepository domain.UserRepository
	ContextTimeout time.Duration
}

func (li *LoginInterceptor) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	c, cancel := context.WithTimeout(ctx, li.ContextTimeout)
	defer cancel()

	return li.UserRepository.GetByEmail(c, email)
}

func (li *LoginInterceptor) CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error) {
	return utils.CreateAccessToken(user, secret, exp)
}

func (li *LoginInterceptor) CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error) {
	return utils.CreateRefreshToken(user, secret, exp)
}
