package models

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID              string `json:"id" bson:"id"`
	Username        string `json:"username" bson:"username"`
	Email           string `json:"email" bson:"email"`
	Password        string `json:"password" bson:"password"`
	ConfirmPassword string `json:"confirmPassword" bson:"-"`
	Approval        bool   `json:"approval" bson:"approval"`
}
