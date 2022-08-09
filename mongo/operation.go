package mongo

import (
	"context"
	"fmt"
	"reflect"

	"mongo-orm/errorType"
	"mongo-orm/util"

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

func (col *Collection) decodeCursor(cursor *mongo.Cursor, t reflect.Type) interface{} {
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 10)
	for cursor.Next(context.Background()) {
		doc := reflect.New(t).Interface()
		if err := cursor.Decode(doc); err != nil {
			fmt.Println(err.Error())
		}
		slice = reflect.Append(slice, reflect.ValueOf(doc).Elem())
	}
	return slice.Interface()
}

func (col *Collection) decodeCursorAll(cursor *mongo.Cursor, slice interface{}) error {
	if err := cursor.All(context.Background(), slice); err != nil {
		return err
	}
	return nil
}

func (col *Collection) findAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return col.Collection.Find(ctx, filter, opts...)
}

func (col *Collection) FindAll(requiredExample interface{}, filter interface{}, opts ...*options.FindOptions) (interface{}, error) {
	cursor, err := col.findAll(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return col.decodeCursor(cursor, util.GetInterfaceType(requiredExample)), nil
}

func (col *Collection) FindAllTS(requiredExample interface{}, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOptions) (interface{}, error) {
	cursor, err := col.findAll(*sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return col.decodeCursor(cursor, util.GetInterfaceType(requiredExample)), nil
}

func (col *Collection) FindAllAndDecode(resultSlice interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cursor, err := col.findAll(context.Background(), filter, opts...)
	if err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if err = col.decodeCursorAll(cursor, resultSlice); err != nil {
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}

	return nil
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
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	if updateResult.ModifiedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotModifiedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.Collection.UpdateMany(context.Background(), filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	if updateResult.ModifiedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotModifiedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection) ReplaceOne(filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	result, err := col.Collection.ReplaceOne(context.Background(), filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	if result.MatchedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, document)
	}
	if result.ModifiedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotModifiedAnyErr, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.Collection.DeleteOne(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.Collection.DeleteMany(context.Background(), filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}
