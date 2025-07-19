package models

type Master struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Type        string `bson:"type" json:"type"`
	Value       string `bson:"value" json:"value"`
	Description string `bson:"description" json:"description"`
	IsActive    bool   `bson:"is_active" json:"is_active"`
}
