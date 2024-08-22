package models

type Blog struct {
	ID        string `json:"id" bson:"_id"`
	Title     string `json:"title" bson:"title"`
	Body      string `json:"body" bson:"body"`
	CreatedBy string `json:"createdBy" bson:"createdBy"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy string `json:"updatedBy" bson:"updatedBy"`
}
