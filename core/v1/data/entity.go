package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Data struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Identifier string             `bson:"identifier" json:"identifier" binding:"required"`
	Url        string             `bson:"url" json:"url" binding:"required,url"`
	UrlSource  string             `bson:"identifier" json:"identifier" binding:"required"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type AssembledScreen struct {
	Body      interface{} `json:"body"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
