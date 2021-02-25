package configuration

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("Configuration not found")

type Repository struct {
	db *mongo.Database
}

func NewRepository(database *mongo.Database) *Repository {
	return &Repository{database}
}

func (r *Repository) col() *mongo.Collection {
	return r.db.Collection("configuration")
}

func (r *Repository) Find(ctx context.Context, in interface{}) (*Configuration, error) {
	var data Configuration
	if err := r.col().FindOne(ctx, in).Decode(&data); err != nil {
		return nil, ErrNotFound
	}

	return &data, nil
}
