package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/utils"
	"time"
)

type SignupInterceptor struct {
	UserRepository domain.UserRepository
	ContextTimeout time.Duration
}

func (si *SignupInterceptor) Store(ctx context.Context, user *domain.User) error {
	c, cancel := context.WithTimeout(ctx, si.ContextTimeout)
	defer cancel()

	return si.UserRepository.Create(c, user)
}

func (si *SignupInterceptor) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	c, cancel := context.WithTimeout(ctx, si.ContextTimeout)
	defer cancel()

	return si.UserRepository.GetByEmail(c, email)
}

func (si *SignupInterceptor) CreateAccessToken(user *domain.User, secret string, exp int) (accessToken string, err error) {
	accessToken, err = utils.CreateAccessToken(user, secret, exp)

	return
}

func (si *SignupInterceptor) CreateRefreshToken(user *domain.User, secret string, exp int) (refreshToken string, err error) {
	refreshToken, err = utils.CreateRefreshToken(user, secret, exp)

	return
}
