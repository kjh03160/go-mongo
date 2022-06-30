package mongo

import (
	"context"

	"mongo-orm/errorType"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	*mongo.Collection
}

func (col *Collection) evaluateAndDecodeSingleResult(result *mongo.SingleResult, v interface{}) error {
	if result == nil {
		return errorType.SingleResultErr
	}
	if err := result.Decode(v); err != nil {
		return err
	}
	return nil
}

func (col *Collection) FindAll(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := col.Collection.Find(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return cursor, nil
}

func (col *Collection) FindOne(data, filter interface{}, opts ...*options.FindOneOptions) error {
	singleResult := col.Collection.FindOne(context.Background(), filter, opts...)
	if err := col.evaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return nil
}

func (col *Collection) FindOneAndModify(data, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	singleResult := col.Collection.FindOneAndUpdate(context.Background(), filter, update, opts...)
	if err := col.evaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	return nil
}

func (col *Collection) FindOneAndReplace(data, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	singleResult := col.Collection.FindOneAndReplace(context.Background(), filter, replacement, opts...)
	if err := col.evaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, replacement, nil)
	}
	return nil
}

func (col *Collection) FindOneAndDelete(data, filter interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	singleResult := col.Collection.FindOneAndDelete(context.Background(), filter, opts...)
	if err := col.evaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return nil
}

func (col *Collection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.Collection.UpdateOne(context.Background(), filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.Collection.UpdateMany(context.Background(), filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection) ReplaceOne(filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	result, err := col.Collection.ReplaceOne(context.Background(), filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.Collection.DeleteOne(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.Collection.DeleteMany(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}
