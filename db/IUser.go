package db

import (
	"blog_app/models"
	"blog_app/mycontext"
)

type IUser interface {
	GetUserByEmail(ctx mycontext.Context, email string) (models.User, error)
	CreateUser(ctx mycontext.Context, user models.User) error
}
