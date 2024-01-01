package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTask = "tasks"
)

type Task struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"-"`
	Title  string             `bson:"title" form:"title" binding:"required" json:"title"`
}

type TaskRepository interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}

type TaskUseCase interface {
	Create(c context.Context, task *Task) error
	FetchByUserID(c context.Context, userID string) ([]Task, error)
}
