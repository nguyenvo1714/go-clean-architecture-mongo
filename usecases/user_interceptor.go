package usecases

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"time"
)

type UserInterceptor struct {
	UserRepository domain.UserRepository
	ContextTimeout time.Duration
}

func (ui *UserInterceptor) Store(ctx context.Context, user *domain.User) error {
	c, cancel := context.WithTimeout(ctx, ui.ContextTimeout)
	defer cancel()

	return ui.UserRepository.Create(c, user)
}

func (ui *UserInterceptor) Fetch(ctx context.Context) ([]domain.User, error) {
	c, cancel := context.WithTimeout(ctx, ui.ContextTimeout)
	defer cancel()

	return ui.UserRepository.Fetch(c)
}

func (ui *UserInterceptor) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	c, cancel := context.WithTimeout(ctx, ui.ContextTimeout)
	defer cancel()

	return ui.UserRepository.GetByEmail(c, email)
}

func (ui *UserInterceptor) GetByID(ctx context.Context, id string) (domain.User, error) {
	c, cancel := context.WithTimeout(ctx, ui.ContextTimeout)
	defer cancel()

	return ui.UserRepository.GetByID(c, id)
}
