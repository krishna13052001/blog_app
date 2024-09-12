package domain

import (
	"blog_app/log"
	"blog_app/models"
	"blog_app/mycontext"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *domainService) CreateBlog(ctx mycontext.Context, blog models.Blog) error {
	blog.ID = uuid.New().String()
	blog.CreatedAt = time.Now().UnixMilli()
	blog.UpdatedAt = time.Now().UnixMilli()
	log.GenericInfo(ctx, "Creating blog", log.FieldsMap{"blog": blog})
	return s.DB.CreateBlog(ctx, blog)
}

func (s *domainService) GetBlog(ctx mycontext.Context, start string) ([]models.Blog, error) {
	blogs, err := s.DB.GetBlog(ctx, start)
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
	if blog.ID == "" {
		return errors.New("invalid blog id")
	}
	blog.UpdatedAt = time.Now().UnixMilli()
	err := s.DB.UpdateBlog(ctx, blog)
	if err != nil {
		log.GenericError(ctx, errors.New("error updating blog"), log.FieldsMap{"error": err.Error()})
		return err
	}
	return nil
}

func (s *domainService) ValidateUser(ctx mycontext.Context, cred models.Credentials) (models.User, error) {
	user, err := s.DB.GetUserByEmail(ctx, cred.Email)
	if err != nil {
		log.GenericError(ctx, errors.New("error getting user"), log.FieldsMap{"error": err.Error()})
		return models.User{}, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		log.GenericError(ctx, errors.New("error comparing password"), log.FieldsMap{"error": err.Error()})
		return models.User{}, errors.New("invalid username or password")
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *domainService) UserExists(ctx mycontext.Context, email string) (bool, error) {
	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		log.GenericError(ctx, errors.New("error getting user"), log.FieldsMap{"error": err.Error()})
		return false, err
	}
	if user.Email == email {
		return true, nil
	}
	return false, nil
}

func (s *domainService) RegisterUser(ctx mycontext.Context, user models.User) error {
	user.ID = uuid.New().String()
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.GenericError(ctx, errors.New("error hashing password"), log.FieldsMap{"error": err.Error()})
		return err
	}
	user.Password = hashedPassword
	err = s.DB.CreateUser(ctx, user)
	if err != nil {
		log.GenericError(ctx, errors.New("error creating user"), log.FieldsMap{"error": err.Error()})
		return err
	}
	return nil
}
