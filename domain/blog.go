package domain

import (
	"blog_app/log"
	"blog_app/models"
	"blog_app/mycontext"
	"errors"
	"github.com/google/uuid"
	"time"
)

func (s *domainService) CreateBlog(ctx mycontext.Context, blog models.Blog) error {
	blog.ID = uuid.New().String()
	blog.UpdatedAt = time.Now().UnixMilli()
	log.GenericInfo(ctx, "Creating blog", log.FieldsMap{"blog": blog})
	return s.DB.CreateBlog(ctx, blog)
}

func (s *domainService) GetBlog(ctx mycontext.Context) ([]models.Blog, error) {
	blogs, err := s.DB.GetBlog(ctx)
	if err != nil {
		log.GenericError(ctx, errors.New("error getting blog"), log.FieldsMap{"error": err.Error()})
		return nil, err
	}
	return blogs, nil
}

func (s *domainService) GetBlogById(ctx mycontext.Context, id string) (models.Blog, error) {
	blog, err := s.DB.GetBlogById(ctx, id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			readErr := errors.New("no blog found with the given ID")
			log.GenericError(ctx, readErr, log.FieldsMap{"context": ctx, "id": id})
			return models.Blog{}, readErr
		}
		log.GenericError(ctx, errors.New("error getting blog"), log.FieldsMap{"error": err.Error()})
		return models.Blog{}, err
	}
	return blog, nil
}

func (s *domainService) DeleteBlog(ctx mycontext.Context, id string) error {
	err := s.DB.DeleteBlog(ctx, id)
	if err != nil {
		log.GenericError(ctx, errors.New("error deleting blog"), log.FieldsMap{"error": err.Error()})
		return err
	}
	return nil
}

func (s *domainService) UpdateBlog(ctx mycontext.Context, blog models.Blog) error {
	blog.UpdatedAt = time.Now().UnixMilli()
	err := s.DB.UpdateBlog(ctx, blog)
	if err != nil {
		log.GenericError(ctx, errors.New("error updating blog"), log.FieldsMap{"error": err.Error()})
		return err
	}
	return nil
}
