package mongo

import (
	"context"

	"mongo-orm/errorType"
	"mongo-orm/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (col *Collection) FindAll(requiredExample interface{}, filter interface{}, opts ...*options.FindOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	cursor, err := col.findAll(ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	result, err := DecodeCursor(cursor, util.GetInterfaceType(requiredExample))
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return result, nil
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
	singleResult := col.findOne(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) FindOneTS(data, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneOptions) error {
	singleResult := col.findOne(*sessCtx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) FindOneAndModify(data, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndModify(ctx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) FindOneAndReplace(data, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndReplace(ctx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) FindOneAndDelete(data, filter interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndDelete(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.insertOne(ctx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.insertMany(ctx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	updateResult, err := col.updateOne(ctx, filter, update, opts...)
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
	updateResult, err := col.updateMany(ctx, filter, update, opts...)
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
	result, err := col.replaceOne(ctx, filter, document, opts...)
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
	deleteResult, err := col.deleteOne(ctx, filter, opts...)
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
	deleteResult, err := col.deleteMany(ctx, filter, opts...)
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
	count, err := col.countDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	count, err := col.estimatedDocumentCount(ctx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}
