package db

import (
	"blog_app/constants"
	"blog_app/log"
	"blog_app/models"
	"blog_app/mycontext"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoServices) CreateBlog(ctx mycontext.Context, blog models.Blog) error {
	_, err := m.Db.CreateOne(ctx, constants.BlogAppDatabase, constants.BlogCollection, blog)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while creating blog"), log.FieldsMap{"blog": blog})
		return err
	}
	return nil
}

func (m *MongoServices) GetBlog(ctx mycontext.Context) ([]models.Blog, error) {
	var blogs []models.Blog
	err := m.Db.ReadAll(ctx, constants.BlogAppDatabase, constants.BlogCollection, nil, &blogs)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while getting blog"), nil)
		return nil, err
	}
	return blogs, nil
}

func (m *MongoServices) GetBlogById(ctx mycontext.Context, id string) (models.Blog, error) {
	objectId, _ := primitive.ObjectIDFromHex(id)
	var blog models.Blog
	err := m.Db.ReadOne(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": objectId}, &blog)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while updating blog"), log.FieldsMap{"blog": blog})
		return models.Blog{}, err
	}
	return blog, nil
}

func (m *MongoServices) DeleteBlog(ctx mycontext.Context, id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := m.Db.Delete(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": objectId})
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while deleting blog"), log.FieldsMap{"id": id})
		return err
	}
	return nil
}

func (m *MongoServices) UpdateBlog(ctx mycontext.Context, blog models.Blog) error {
	if blog.ID == "" {
		return errors.New("invalid blog id")
	}
	objectId, _ := primitive.ObjectIDFromHex(blog.ID)
	updateOptions := options.Update()
	updateOptions.SetUpsert(false)
	updateMap := bson.M{
		"title":     blog.Title,
		"body":      blog.Body,
		"batch":     blog.Batch,
		"jobType":   blog.JobType,
		"payRange":  blog.PayRange,
		"applyLink": blog.ApplyLink,
		"updatedAt": blog.UpdatedAt,
	}
	_, err := m.Db.Update(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": objectId}, bson.M{"$set": updateMap}, updateOptions)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while updating blog"), log.FieldsMap{"blog": blog})
		return err
	}
	return nil
}
