package mongo

import (
	"context"

	"mongo-orm/errorType"
	"mongo-orm/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collection struct {
	*mongo.Collection
}

func (col *Collection) FindAll(requiredExample interface{}, filter interface{}, opts ...*options.FindOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	cursor, err := col.findAll(ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return DecodeCursor(cursor, util.GetInterfaceType(requiredExample)), nil
}

func (col *Collection) FindAllTS(requiredExample interface{}, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOptions) (interface{}, error) {
	cursor, err := col.findAll(*sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return DecodeCursor(cursor, util.GetInterfaceType(requiredExample)), nil
}

func (col *Collection) FindAllAndDecode(resultSlice interface{}, filter interface{}, opts ...*options.FindOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	cursor, err := col.findAll(ctx, filter, opts...)
	if err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if err = DecodeCursorAll(cursor, resultSlice); err != nil {
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) FindOne(data, filter interface{}, opts ...*options.FindOneOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.Collection.FindOne(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return nil
}

func (col *Collection) FindOneAndModify(data, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	return nil
}

func (col *Collection) FindOneAndReplace(data, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, replacement, nil)
	}
	return nil
}

func (col *Collection) FindOneAndDelete(data, filter interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.Collection.FindOneAndDelete(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return nil
}

func (col *Collection) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.Collection.InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.Collection.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	updateResult, err := col.Collection.UpdateOne(ctx, filter, update, opts...)
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	updateResult, err := col.Collection.UpdateMany(ctx, filter, update, opts...)
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	result, err := col.Collection.ReplaceOne(ctx, filter, document, opts...)
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	deleteResult, err := col.Collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	deleteResult, err := col.Collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	count, err := col.Collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	count, err := col.Collection.EstimatedDocumentCount(ctx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}

func (col *Collection) findAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return col.Collection.Find(ctx, filter, opts...)
}

func Transaction(client *mongo.Client, function func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	opt := options.TransactionOptions{}
	opt.SetReadPreference(readpref.Primary())
	_, err = session.WithTransaction(context.Background(), function, &opt)
	if err != nil {
		return err
	}
	return nil
}
