package db

import (
	"blog_app/connectionMgr"
	"reflect"
	"strings"
)

var _ Service = (*MongoServices)(nil)

type MongoServices struct {
	Db connectionMgr.MongoDB
}

func NewMongoService(mongoClient connectionMgr.MongoDB) Service {
	return &MongoServices{
		Db: mongoClient,
	}
}

func GetBsonTag(input interface{}, jsonTag string) string {
	t := reflect.TypeOf(input)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := strings.Split(f.Tag.Get("json"), ",")[0]
		if tag == jsonTag {
			return strings.Split(f.Tag.Get("bson"), ",")[0]
		}
	}
	return jsonTag
}
