package mongo

import (
	"context"

	"mongo-orm/errorType"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (col *Collection[T]) FindAll(filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	cursor, err := col.findAll(ctx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindOne(data, filter interface{}, opts ...*options.FindOneOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOne(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndModify(data, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndModify(ctx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndReplace(data, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndReplace(ctx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndDelete(data, filter interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndDelete(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.insertOne(ctx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection[T]) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (interface{}, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	insertOneResult, err := col.insertMany(ctx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection[T]) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	updateResult, err := col.updateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	updateResult, err := col.updateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) ReplaceOne(filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	result, err := col.replaceOne(ctx, filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	if result.MatchedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection[T]) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
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

func (col *Collection[T]) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
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

func (col *Collection[T]) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	count, err := col.countDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	count, err := col.estimatedDocumentCount(ctx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	bulkWriteResult, err := col.bulkWrite(ctx, models, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, models, nil)
	}
	return bulkWriteResult, nil
}

func (col *Collection[T]) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	cursor, err := col.aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), pipeline, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), pipeline, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindAllWithTrx(filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := col.findAll(*sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return resultSlice, nil
}

func (col *Collection[T]) FindOneWithTrx(data, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneOptions) error {
	singleResult := col.findOne(*sessCtx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Unwrap(err) == context.DeadlineExceeded {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndModifyWithTrx(data, filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndUpdateOptions) error {
	singleResult := col.findOneAndModify(*sessCtx, filter, update, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndReplaceWithTrx(data, filter interface{}, replacement interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndReplaceOptions) error {
	singleResult := col.findOneAndReplace(*sessCtx, filter, replacement, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) FindOneAndDeleteWithTrx(data, filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.FindOneAndDeleteOptions) error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer ctxCancel()
	singleResult := col.findOneAndDelete(ctx, filter, opts...)
	if err := EvaluateAndDecodeSingleResult(singleResult, data); err != nil {
		if err == mongo.ErrNoDocuments || errors.Is(err, context.DeadlineExceeded) {
			return errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
		}
		return errorType.DecodeError(col.Name(), filter, nil, nil, err)
	}
	return nil
}

func (col *Collection[T]) InsertOneWithTrx(document interface{}, sessCtx *mongo.SessionContext, opts ...*options.InsertOneOptions) (interface{}, error) {
	insertOneResult, err := col.insertOne(*sessCtx, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, document)
	}
	return insertOneResult.InsertedID, nil
}

func (col *Collection[T]) InsertManyWithTrx(documents []interface{}, sessCtx *mongo.SessionContext, opts ...*options.InsertManyOptions) (interface{}, error) {
	insertOneResult, err := col.insertMany(*sessCtx, documents, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, documents)
	}
	return insertOneResult.InsertedIDs, nil
}

func (col *Collection[T]) UpdateOneWithTrx(filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.updateOne(*sessCtx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) UpdateManyWithTrx(filter interface{}, update interface{}, sessCtx *mongo.SessionContext, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	updateResult, err := col.updateMany(*sessCtx, filter, update, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, update, nil)
	}
	if updateResult.MatchedCount == 0 {
		return updateResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, update, nil)
	}
	return updateResult, nil
}

func (col *Collection[T]) ReplaceOneWithTrx(filter interface{}, document interface{}, sessCtx *mongo.SessionContext, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	result, err := col.replaceOne(*sessCtx, filter, document, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, document)
	}
	if result.MatchedCount == 0 {
		return result, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, document)
	}
	return result, nil
}

func (col *Collection[T]) DeleteOneWithTrx(filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.deleteOne(*sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) DeleteManyWithTrx(filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	deleteResult, err := col.deleteMany(*sessCtx, filter, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	if deleteResult.DeletedCount == 0 {
		return deleteResult, errorType.ParseAndReturnDBError(errorType.NotMatchedAnyErr, col.Name(), filter, nil, nil)
	}
	return deleteResult, nil
}

func (col *Collection[T]) CountDocumentsWithTrx(filter interface{}, sessCtx *mongo.SessionContext, opts ...*options.CountOptions) (int, error) {
	count, err := col.countDocuments(*sessCtx, filter, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), filter, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) EstimatedDocumentCountWithTrx(sessCtx *mongo.SessionContext, opts ...*options.EstimatedDocumentCountOptions) (int, error) {
	count, err := col.estimatedDocumentCount(*sessCtx, opts...)
	if err != nil {
		return 0, errorType.ParseAndReturnDBError(err, col.Name(), nil, nil, nil)
	}
	return int(count), nil
}

func (col *Collection[T]) BulkWriteWithTrx(models []mongo.WriteModel, sessCtx *mongo.SessionContext, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	bulkWriteResult, err := col.bulkWrite(*sessCtx, models, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), nil, models, nil)
	}
	return bulkWriteResult, nil
}

func (col *Collection[T]) AggregateWithTrx(pipeline interface{}, sessCtx *mongo.SessionContext, opts ...*options.AggregateOptions) ([]T, error) {
	cursor, err := col.aggregate(*sessCtx, pipeline, opts...)
	if err != nil {
		return nil, errorType.ParseAndReturnDBError(err, col.Name(), pipeline, nil, nil)
	}
	resultSlice, err := DecodeCursor[T](cursor)
	if err != nil {
		return nil, errorType.DecodeError(col.Name(), pipeline, nil, nil, err)
	}
	return resultSlice, nil
}
