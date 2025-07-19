package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Document struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Size       int64              `bson:"size" json:"size"`
	Type       string             `bson:"type" json:"type"`
	FolderID   primitive.ObjectID `bson:"folder_id" json:"folder_id"`
	UploadedAt time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	URL        string             `bson:"url" json:"url"`
}
