package data

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrNotFound = errors.New("data not found")

type Repository struct {
	db *mongo.Database
}

func NewRepository(database *mongo.Database) *Repository {
	return &Repository{database}
}

func (r *Repository) col() *mongo.Collection {
	return r.db.Collection("data")
}

func (r *Repository) Find(ctx context.Context, in interface{}) (*Data, error) {
	var data Data
	if err := r.col().FindOne(ctx, in).Decode(&data); err != nil {
		return nil, ErrNotFound
	}

	return &data, nil
}
