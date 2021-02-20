package configuration

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type KeyValue struct {
	Index string `bson:"index" json:"index"`
	Value string `bson:"value" json:"value"`
}

type ServiceRequest struct {
	Name        string     `bson:"name" json:"name" binding:"required"`
	Url         string     `bson:"url" json:"url" binding:"required,url"`
	Body        string     `bson:"body,omitempty" json:"body"`
	Headers     []KeyValue `bson:"headers,omitempty" json:"headers"`
	Mapping     []KeyValue `bson:"mapping,omitempty" json:"mapping"`
	Validation  []KeyValue `bson:"validation,omitempty" json:"validation"`
	Method      string     `bson:"method" json:"method" binding:"required"`
	ContentType string     `bson:"content_type" json:"content_type"`
}

type Data struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Identifier string             `bson:"identifier" json:"identifier" binding:"required"`
	Source     ServiceRequest     `bson:"source" json:"source" binding:"required"`
	Crawler    []ServiceRequest   `bson:"crawler" json:"crawler" binding:"required"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type AssembledScreen struct {
	Identifier string      `json:"identifier"`
	Data       interface{} `json:"data"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func (crawler ServiceRequest) HasMapping() bool {
	return crawler.Mapping != nil
}
