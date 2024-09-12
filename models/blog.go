package models

type Blog struct {
	ApplyLink   string `json:"applyLink" bson:"applyLink"`
	Body        string `json:"body" bson:"body"`
	Batch       string `json:"batch" bson:"batch"`
	Company     string `json:"company" bson:"company"`
	CreatedAt   int64  `json:"createdAt" bson:"createdAt"`
	ID          string `json:"id" bson:"_id"`
	JobCategory string `json:"jobCategory" json:"jobCategory"`
	JobType     string `json:"jobType" bson:"jobType"`
	Location    string `json:"location" bson:"location"`
	PayRange    string `json:"payRange" bson:"payRange"`
	Title       string `json:"title" bson:"title"`
	UpdatedAt   int64  `json:"updatedAt" bson:"updatedAt"`
}
