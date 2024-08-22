package connectionMgr

import (
	"blog_app/log"
	"blog_app/mycontext"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (client *mongoService) ReadOne(ctx mycontext.Context, database, collectionName string, filter, data interface{}) error {
	collection := client.DB.Database(database).Collection(collectionName)
	err := collection.FindOne(ctx.Context, filter).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (client *mongoService) ReadAll(ctx mycontext.Context, database, collectionName string, filter, data interface{}, opts ...*options.FindOptions) error {
	var findOptions *options.FindOptions
	if len(opts) > 0 {
		findOptions = opts[0]
	}

	if filter == nil {
		filter = bson.M{}
	}

	collection := client.DB.Database(database).Collection(collectionName)
	cursor, err := collection.Find(ctx.Context, filter, findOptions)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx.Context)

	err = cursor.All(ctx.Context, data)
	if err != nil {
		return err
	}

	return nil
}

func (client *mongoService) CreateOne(ctx mycontext.Context, database, collectionName string, d interface{}) (*mongo.InsertOneResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	insertOneRslt, err := collection.InsertOne(ctx.Context, d)
	if err != nil {
		return nil, err
	}

	return insertOneRslt, nil
}

func (client *mongoService) CreateMany(ctx mycontext.Context, database, collectionName string, d []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	insertManyRslt, err := collection.InsertMany(ctx.Context, d)
	if err != nil {
		return nil, err
	}
	return insertManyRslt, nil
}

func (client *mongoService) Update(ctx mycontext.Context, database, collectionName string, filter, update interface{}, options ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	updateResult, err := collection.UpdateOne(ctx.Context, filter, update, options...)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

func (client *mongoService) Upsert(ctx mycontext.Context, database, collectionName string, filter,
	update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	updateOptions := options.Update()
	if len(opts) >= 1 {
		updateOptions = opts[0]
	}
	updateOptions.SetUpsert(true)
	updateResult, err := collection.UpdateOne(ctx.Context, filter, update, updateOptions)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

func (client *mongoService) UpdateAndReturn(ctx mycontext.Context, database, collectionName string, filter, update, data interface{}) error {
	collection := client.DB.Database(database).Collection(collectionName)
	after := options.After

	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	err := collection.FindOneAndUpdate(ctx.Context, filter, update, &opts).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func (client *mongoService) UpdateAll(ctx mycontext.Context, database, collectionName string, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	updateResult, err := collection.UpdateMany(ctx.Context, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

func (client *mongoService) FindOneAndUpdate(ctx mycontext.Context, database, collectionName string, filter, update, data interface{},
	opts ...*options.FindOneAndUpdateOptions) error {
	after := options.After
	option := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	if len(opts) > 0 {
		option = opts[0]
	}
	collection := client.DB.Database(database).Collection(collectionName)
	result := collection.FindOneAndUpdate(ctx.Context, filter, update, option)
	if result.Err() != nil {
		return result.Err()
	}

	decodeErr := result.Decode(data)
	if decodeErr != nil {
		return decodeErr
	}
	return nil
}

func (client *mongoService) Delete(ctx mycontext.Context, database, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	deleteResult, err := collection.DeleteOne(ctx.Context, filter)

	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

func (client *mongoService) DeleteAll(ctx mycontext.Context, database, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	deleteResult, err := collection.DeleteMany(ctx.Context, filter)
	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

func (client *mongoService) CountDocuments(ctx mycontext.Context, database, collectionName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var countOptions *options.CountOptions
	if len(opts) > 0 {
		countOptions = opts[0]
	}

	collection := client.DB.Database(database).Collection(collectionName)
	count, err := collection.CountDocuments(ctx.Context, filter, countOptions)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (client *mongoService) BulkWrite(ctx mycontext.Context, database, collectionName string, operations []mongo.WriteModel, bulkOption *options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	var err error
	collection := client.DB.Database(database).Collection(collectionName)
	result, err := collection.BulkWrite(ctx.Context, operations, bulkOption)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (client *mongoService) Exist(ctx mycontext.Context, database, collectionName string, filter interface{}) (bool, error) {

	var i interface{}
	collection := client.DB.Database(database).Collection(collectionName)
	err := collection.FindOne(ctx.Context, filter).Decode(&i)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (client *mongoService) GetDistinct(ctx mycontext.Context, database, collectionName, fieldName string, filter interface{}) (interface{}, error) {
	collection := client.DB.Database(database).Collection(collectionName)
	result, err := collection.Distinct(ctx.Context, fieldName, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client *mongoService) AggregateAll(ctx mycontext.Context, database, collectionName string, query, data interface{}, options ...*options.AggregateOptions) error {
	collection := client.DB.Database(database).Collection(collectionName)
	cursor, err := collection.Aggregate(ctx.Context, query, options...)
	if err != nil {
		return err
	}
	err = cursor.All(ctx.Context, data)
	return err
}

func (client *mongoService) CreateIndexes(ctx mycontext.Context, database, collectionName string, operations []mongo.IndexModel) ([]string, error) {
	var err error
	collection := client.DB.Database(database).Collection(collectionName)
	result, err := collection.Indexes().CreateMany(ctx.Context, operations)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "CreateIndexes: Error in creating index"), log.FieldsMap{"collectionName": collectionName})
		return nil, err
	}
	return result, nil
}

func (client *mongoService) CreateIndex(ctx mycontext.Context, database, collectionName string, operations mongo.IndexModel) (string, error) {
	var err error
	collection := client.DB.Database(database).Collection(collectionName)
	result, err := collection.Indexes().CreateOne(ctx.Context, operations)
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "CreateIndex: Error in creating index"), log.FieldsMap{"collectionName": collectionName})
		return "", err
	}
	return result, nil
}

func (client *mongoService) Disconnect(ctx mycontext.Context) error {
	return client.DB.Disconnect(ctx)
}
