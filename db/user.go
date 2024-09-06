package db

import (
	"blog_app/constants"
	"blog_app/log"
	"blog_app/models"
	"blog_app/mycontext"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoServices) GetUserByEmail(ctx mycontext.Context, email string) (models.User, error) {
	var user models.User
	err := m.Db.ReadOne(ctx, constants.BlogAppDatabase, constants.UserCollection, bson.M{"email": email}, &user)
	if err != nil {
		return models.User{}, err
	}
	if user.Email == "" {
		return models.User{}, errors.New("Invalid user details")
	}
	return user, nil
}

func (m *MongoServices) CreateUser(ctx mycontext.Context, user models.User) error {
	_, err := m.Db.CreateOne(ctx, constants.BlogAppDatabase, constants.UserCollection, user)
	if err != nil {
		log.GenericError(ctx, errors.New("error creating user"+err.Error()), log.FieldsMap{"user": user})
		return err
	}
	return nil
}
