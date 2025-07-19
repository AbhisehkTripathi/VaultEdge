package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Folder struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name         string              `bson:"name" json:"name"`
	ParentID     *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	CreatedAt    time.Time           `bson:"created_at" json:"created_at"`
	DocumentCount int                `bson:"document_count" json:"document_count"`
}
