package db

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
