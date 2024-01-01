package repositories

import (
	"context"
	"go-learning/clean-architecture-mongo/database"
	"go-learning/clean-architecture-mongo/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	database   database.Database
	collection string
}

func NewUserRepository(database database.Database, collection string) *UserRepository {
	return &UserRepository{
		database:   database,
		collection: collection,
	}
}

func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := ur.database.Collection(ur.collection).InsertOne(ctx, user)

	return err
}

func (ur *UserRepository) Fetch(ctx context.Context) ([]domain.User, error) {
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := ur.database.Collection(ur.collection).Find(ctx, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User
	err = cursor.All(ctx, &users)

	if err != nil || users == nil {
		return []domain.User{}, err
	}

	return users, err
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	result := ur.database.Collection(ur.collection).FindOne(ctx, bson.M{"email": email})
	err := result.Decode(&user)

	return user, err
}

func (ur *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	idHex, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return domain.User{}, err
	}

	result := ur.database.Collection(ur.collection).FindOne(ctx, bson.M{"_id": idHex})
	err = result.Decode(&user)

	return user, err
}
