package repositories

import (
	"context"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository struct {
	database   database.Database
	collection string
}

func NewTaskRepository(database database.Database, collection string) *TaskRepository {
	return &TaskRepository{
		database:   database,
		collection: collection,
	}
}

func (tr *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := tr.database.Collection(tr.collection).InsertOne(ctx, task)

	return err
}

func (tr *TaskRepository) FetchByUserID(ctx context.Context, userID string) ([]domain.Task, error) {
	var tasks []domain.Task
	idHex, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return tasks, err
	}

	cursor, err := tr.database.Collection(tr.collection).Find(ctx, bson.M{"user_id": idHex})

	if err != nil {
		return []domain.Task{}, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil || tasks == nil {
		return []domain.Task{}, err
	}

	return tasks, err
}
