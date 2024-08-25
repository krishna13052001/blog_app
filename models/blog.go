package models

type Blog struct {
	ID        string `json:"id" bson:"_id"`
	Title     string `json:"title" bson:"title"`
	Body      string `json:"body" bson:"body"`
	Batch     string `json:"batch" bson:"batch"`
	JobType   string `json:"jobType" bson:"jobType"`
	Location  string `json:"location" bson:"location"`
	PayRange  string `json:"payRange" bson:"payRange"`
	ApplyLink string `json:"applyLink" bson:"applyLink"`
	UpdatedAt int64  `json:"updatedAt" bson:"updatedAt"`
}
