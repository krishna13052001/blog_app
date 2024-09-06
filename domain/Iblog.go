package domain

import (
	"blog_app/models"
	"blog_app/mycontext"
)

type IBlog interface {
	CreateBlog(ctx mycontext.Context, blog models.Blog) error
	GetBlog(ctx mycontext.Context, start string) ([]models.Blog, error)
	GetBlogById(ctx mycontext.Context, id string) (models.Blog, error)
	DeleteBlog(ctx mycontext.Context, id string) error
	UpdateBlog(ctx mycontext.Context, blog models.Blog) error
}

type IUser interface {
	ValidateUser(ctx mycontext.Context, cred models.Credentials) (models.User, error)
	UserExists(ctx mycontext.Context, email string) (bool, error)
	RegisterUser(ctx mycontext.Context, user models.User) error
}
