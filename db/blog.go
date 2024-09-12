package db

import (
	"blog_app/constants"
	"blog_app/log"
	"blog_app/models"
	"blog_app/mycontext"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

func (m *MongoServices) CreateBlog(ctx mycontext.Context, blog models.Blog) error {
	_, err := m.Db.CreateOne(ctx, constants.BlogAppDatabase, constants.BlogCollection, blog)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while creating blog"), log.FieldsMap{"blog": blog})
		return err
	}
	return nil
}

func (m *MongoServices) GetBlog(ctx mycontext.Context, start string) ([]models.Blog, error) {
	var blogs []models.Blog
	findOptions := options.Find()
	findOptions.SetLimit(15)
	findOptions.SetSort(bson.D{{"createdAt", -1}})
	skipLimit, err := strconv.ParseInt(start, 10, 64)
	if start != "" && err == nil {
		findOptions.SetSkip(skipLimit)
	}
	err = m.Db.ReadAll(ctx, constants.BlogAppDatabase, constants.BlogCollection, nil, &blogs, findOptions)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while getting blog"), nil)
		return nil, err
	}
	return blogs, nil
}

func (m *MongoServices) GetBlogById(ctx mycontext.Context, id string) (models.Blog, error) {
	var blog models.Blog
	err := m.Db.ReadOne(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": id}, &blog)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while updating blog"), log.FieldsMap{"blog": blog})
		return models.Blog{}, err
	}
	return blog, nil
}

func (m *MongoServices) DeleteBlog(ctx mycontext.Context, id string) error {
	_, err := m.Db.Delete(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": id})
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while deleting blog"), log.FieldsMap{"id": id})
		return err
	}
	return nil
}

func (m *MongoServices) UpdateBlog(ctx mycontext.Context, blog models.Blog) error {
	updateOptions := options.Update()
	updateOptions.SetUpsert(false)
	updateMap := bson.M{
		"title":       blog.Title,
		"body":        blog.Body,
		"batch":       blog.Batch,
		"jobType":     blog.JobType,
		"payRange":    blog.PayRange,
		"applyLink":   blog.ApplyLink,
		"updatedAt":   blog.UpdatedAt,
		"jobCategory": blog.JobCategory,
		"company":     blog.Company,
	}
	result, err := m.Db.Update(ctx, constants.BlogAppDatabase, constants.BlogCollection, bson.M{"_id": blog.ID}, bson.M{"$set": updateMap}, updateOptions)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "Error while updating blog"), log.FieldsMap{"blog": blog})
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("no blog found")
	} else if result.ModifiedCount == 0 {
		return errors.New("no changes made")
	}
	return nil
}
