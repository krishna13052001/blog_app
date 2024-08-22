package connectionMgr

import (
	"blog_app/log"
	"blog_app/mycontext"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoService struct {
	DB *mongo.Client
}

type MongoDB interface {
	ReadOne(ctx mycontext.Context, database, collectionName string, filter, data interface{}) error
	ReadAll(ctx mycontext.Context, database, collectionName string, filter, data interface{}, opts ...*options.FindOptions) error
	CreateOne(ctx mycontext.Context, database, collectionName string, d interface{}) (*mongo.InsertOneResult, error)
	CreateMany(ctx mycontext.Context, database, collectionName string, d []interface{}) (*mongo.InsertManyResult, error)
	Update(ctx mycontext.Context, database, collectionName string, filter, update interface{}, options ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Upsert(ctx mycontext.Context, database, collectionName string, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateAndReturn(ctx mycontext.Context, database, collectionName string, filter, update, data interface{}) error
	UpdateAll(ctx mycontext.Context, database, collectionName string, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	FindOneAndUpdate(ctx mycontext.Context, database, collectionName string, filter, update, data interface{}, opts ...*options.FindOneAndUpdateOptions) error
	Delete(ctx mycontext.Context, database, collectionName string, filter interface{}) (*mongo.DeleteResult, error)
	DeleteAll(ctx mycontext.Context, database, collectionName string, filter interface{}) (*mongo.DeleteResult, error)
	BulkWrite(ctx mycontext.Context, database, collectionName string, operations []mongo.WriteModel, bulkOption *options.BulkWriteOptions) (*mongo.BulkWriteResult, error)

	// Aggregate functions
	CountDocuments(ctx mycontext.Context, database, collectionName string, filter interface{}, opts ...*options.CountOptions) (int64, error)
	Exist(ctx mycontext.Context, database, collectionName string, filter interface{}) (bool, error)
	GetDistinct(ctx mycontext.Context, database, collectionName, fieldName string, filter interface{}) (interface{}, error)
	AggregateAll(ctx mycontext.Context, database, collectionName string, query, data interface{}, options ...*options.AggregateOptions) error

	CreateIndexes(ctx mycontext.Context, database, collectionName string, operations []mongo.IndexModel) ([]string, error)
	CreateIndex(ctx mycontext.Context, database, collectionName string, operations mongo.IndexModel) (string, error)

	Disconnect(ctx mycontext.Context) error
}

func NewMongoClient(url, appName string, connectOptions map[string]interface{}) (MongoDB, error) {
	ctx := mycontext.NewContext()
	client, err := connect(url, appName, connectOptions)
	if err != nil {
		log.GenericError(ctx, errors.New("can't connect to db: "+err.Error()), nil)
		return &mongoService{DB: client}, err
	}
	return &mongoService{DB: client}, nil
}

func connect(host, appName string, connectOptions map[string]interface{}) (*mongo.Client, error) {
	var client *mongo.Client

	clientOptions := getConnectOption(connectOptions, appName)

	clientOptions = clientOptions.ApplyURI(host)

	err := clientOptions.Validate()
	if err != nil {
		return client, err
	}

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return client, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return client, err
	}

	return client, nil
}

func getConnectOption(connectOptions map[string]interface{}, appName string) *options.ClientOptions {
	servSelectTimeout := time.Duration(15) * time.Second
	connTimeout := time.Duration(10) * time.Second
	idleTime := time.Duration(2) * time.Minute
	socketTimeout := time.Duration(2) * time.Minute
	maxPooling := uint64(100)
	readPrefTagSet := make([]tag.Set, 0)

	for optionKey, val := range connectOptions {
		switch optionKey {
		case "ServerSelectionTimeout":
			if valInt, ok := val.(int64); ok {
				servSelectTimeout = time.Duration(valInt) * time.Second
			}
		case "ConnectTimeout":
			if valInt, ok := val.(int64); ok {
				connTimeout = time.Duration(valInt) * time.Second
			}
		case "MaxConnIdleTime":
			if valInt, ok := val.(int64); ok {
				idleTime = time.Duration(valInt) * time.Second
			}
		case "SocketTimeout":
			if valInt, ok := val.(int64); ok {
				socketTimeout = time.Duration(valInt) * time.Second
			}
		case "MaxPoolSize":
			if valUint64, ok := val.(uint64); ok {
				maxPooling = valUint64
			}
		case "ReadPref":
			if valTagSet, ok := val.([]tag.Set); ok {
				readPrefTagSet = valTagSet
			}
		}
	}

	readPref := getReadPreference(readPrefTagSet)

	return &options.ClientOptions{
		AppName:                &appName,
		ServerSelectionTimeout: &servSelectTimeout,
		ConnectTimeout:         &connTimeout,
		MaxConnIdleTime:        &idleTime,
		MaxPoolSize:            &maxPooling,
		SocketTimeout:          &socketTimeout,
		ReadPreference:         readPref,
	}
}

func getReadPreference(readPrefTagSet []tag.Set) *readpref.ReadPref {
	if len(readPrefTagSet) > 0 {
		readPrefOpts := readpref.WithTagSets(readPrefTagSet...)
		return readpref.SecondaryPreferred(readPrefOpts)
	}
	return nil
}
