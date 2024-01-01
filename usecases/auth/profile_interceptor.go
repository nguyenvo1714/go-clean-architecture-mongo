package auth

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/domain/auth"
	"time"
)

type ProfileInterceptor struct {
	UserRepository domain.UserRepository
	ContextTimeout time.Duration
}

func (pi *ProfileInterceptor) GetProfileByID(ctx context.Context, id string) (*auth.Profile, error) {
	c, cancel := context.WithTimeout(ctx, pi.ContextTimeout)
	defer cancel()

	user, err := pi.UserRepository.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	return &auth.Profile{Email: user.Email, Name: user.Name}, nil
}
