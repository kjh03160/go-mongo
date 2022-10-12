package mongo

import (
	"context"

	"mongo-orm/errorType"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (col *Collection[T]) FindAll(logger Logger, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	cursor, err := col.findAll(logger, ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindOne(logger Logger, data, filter interface{}, opts ...*options.FindOneOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	singleResult := col.findOne(logger, ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndModify(logger Logger, data, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	singleResult := col.findOneAndModify(logger, ctx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndReplace(logger Logger, data, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	singleResult := col.findOneAndReplace(logger, ctx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndDelete(logger Logger, data, filter interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	singleResult := col.findOneAndDelete(logger, ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) InsertOne(logger Logger, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	insertOneResult, err := col.insertOne(logger, ctx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection[T]) InsertMany(logger Logger, documents []interface{}, opts ...*options.InsertManyOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	insertOneResult, err := col.insertMany(logger, ctx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection[T]) UpdateOne(logger Logger, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	updateResult, err := col.updateOne(logger, ctx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) UpdateMany(logger Logger, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	updateResult, err := col.updateMany(logger, ctx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) ReplaceOne(logger Logger, filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	result, err := col.replaceOne(logger, ctx, filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	if result.MatchedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection[T]) DeleteOne(logger Logger, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	deleteResult, err := col.deleteOne(logger, ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) DeleteMany(logger Logger, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	deleteResult, err := col.deleteMany(logger, ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) CountDocuments(logger Logger, filter interface{}, opts ...*options.CountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	count, err := col.countDocuments(logger, ctx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) EstimatedDocumentCount(logger Logger, opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	count, err := col.estimatedDocumentCount(logger, ctx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) BulkWrite(logger Logger, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	bulkWriteResult, err := col.bulkWrite(logger, ctx, models, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, models, nil)
	}
	return bulkWriteResult, nil
}

func (col *Collection[T]) Aggregate(logger Logger, pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_timeout)
	defer ctxCancel()
	cursor, err := col.aggregate(logger, ctx, pipeline, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), pipeline, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), pipeline, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindAllWithTrx(logger Logger, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := col.findAll(logger, *sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindOneWithTrx(logger Logger, data, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneOptions) error {
	singleResult := col.findOne(logger, *sessCtx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Unwrap(err) == context.DeadlineExceeded {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndModifyWithTrx(logger Logger, data, filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndUpdateOptions) error {
	singleResult := col.findOneAndModify(logger, *sessCtx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndReplaceWithTrx(logger Logger, data, filter interface{}, replacement interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndReplaceOptions) error {
	singleResult := col.findOneAndReplace(logger, *sessCtx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndDeleteWithTrx(logger Logger, data, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndDeleteOptions) error {
	singleResult := col.findOneAndDelete(logger, *sessCtx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) InsertOneWithTrx(logger Logger, document interface{}, sessCtx *mongo.SessionContext, opts ...*options.InsertOneOptions) (interface{}, error) {
	insertOneResult, err := col.insertOne(logger, *sessCtx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection[T]) InsertManyWithTrx(logger Logger, documents []interface{}, sessCtx *mongo.SessionContext, opts ...*options.InsertManyOptions) (interface{}, error) {
	insertOneResult, err := col.insertMany(logger, *sessCtx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection[T]) UpdateOneWithTrx(logger Logger, filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.updateOne(logger, *sessCtx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) UpdateManyWithTrx(logger Logger, filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.updateMany(logger, *sessCtx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) ReplaceOneWithTrx(logger Logger, filter interface{}, document interface{}, sessCtx *mongo.SessionContext, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	result, err := col.replaceOne(logger, *sessCtx, filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	if result.MatchedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection[T]) DeleteOneWithTrx(logger Logger, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.deleteOne(logger, *sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) DeleteManyWithTrx(logger Logger, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.deleteMany(logger, *sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) CountDocumentsWithTrx(logger Logger, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.CountOptions) (int, error) {
	count, err := col.countDocuments(logger, *sessCtx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) EstimatedDocumentCountWithTrx(logger Logger, sessCtx *mongo.SessionContext, opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	count, err := col.estimatedDocumentCount(logger, *sessCtx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) BulkWriteWithTrx(logger Logger, models []mongo.WriteModel, sessCtx *mongo.SessionContext, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	bulkWriteResult, err := col.bulkWrite(logger, *sessCtx, models, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, models, nil)
	}
	return bulkWriteResult, nil
}

func (col *Collection[T]) AggregateWithTrx(logger Logger, pipeline interface{}, sessCtx *mongo.SessionContext, opts ...*options.AggregateOptions) ([]T, error) {
	cursor, err := col.aggregate(logger, *sessCtx, pipeline, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), pipeline, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), pipeline, nil, nil, err)
	}
	return resultSlice, nil
}
