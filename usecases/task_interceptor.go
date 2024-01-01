package usecases

import (
	"context"
	"go-learning/clean-architecture-mongo/domain"
	"time"
)

type TaskInterceptor struct {
	TaskRepository domain.TaskRepository
	ContextTimeout time.Duration
}

func (ti *TaskInterceptor) Create(ctx context.Context, task *domain.Task) error {
	c, cancel := context.WithTimeout(ctx, ti.ContextTimeout)
	defer cancel()

	return ti.TaskRepository.Create(c, task)
}

func (ti *TaskInterceptor) FetchByUserID(ctx context.Context, userID string) ([]domain.Task, error) {
	c, cancel := context.WithTimeout(ctx, ti.ContextTimeout)
	defer cancel()

	return ti.TaskRepository.FetchByUserID(c, userID)
}
