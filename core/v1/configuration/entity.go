package configuration

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Data struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Identifier  string             `bson:"identifier" json:"identifier" binding:"required"`
	Url         string             `bson:"url" json:"url" binding:"required,url"`
	UrlSource   string             `bson:"url_source" json:"url_source" binding:"required"`
	DocumentKey string             `bson:"document_key" json:"document_key" binding:"required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type AssembledScreen struct {
	Identifier string      `json:"identifier"`
	Data       interface{} `json:"data"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type Midgard struct {
	Id   string
	Type string
}

// Declare interface
type MidgardRow interface {
	getId() string
}

// get_price function for Courseprice
func (a Midgard) getId() string {
	return a.Id
}
