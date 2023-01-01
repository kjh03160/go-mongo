package wrapper

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (col *Collection[T]) findOne(logger Logger, ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOne(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s findOne slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return singleResult
}

func (col *Collection[T]) findAll(logger Logger, ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	startTime := time.Now()
	cursor, err := col.Collection.Find(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s findAll slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return cursor, err
}

func (col *Collection[T]) findOneAndModify(logger Logger, ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s findOneAndModify slow query(%v) detected. filter: %+v, update: %+v", col.Name(), time.Since(startTime), filter, update))
	}
	return singleResult
}

func (col *Collection[T]) findOneAndReplace(logger Logger, ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s findOneAndReplace slow query(%v) detected. filter: %+v, replacement: %+v", col.Name(), time.Since(startTime), filter, replacement))
	}
	return singleResult
}

func (col *Collection[T]) findOneAndDelete(logger Logger, ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndDelete(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s findOneAndDelete slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return singleResult
}

func (col *Collection[T]) insertOne(logger Logger, ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	startTime := time.Now()
	insertOneResult, err := col.Collection.InsertOne(ctx, document, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s insertOne slow query(%v) detected. document: %+v", col.Name(), time.Since(startTime), document))
	}
	return insertOneResult, err
}

func (col *Collection[T]) insertMany(logger Logger, ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	startTime := time.Now()
	insertOneResult, err := col.Collection.InsertMany(ctx, documents, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s insertMany slow query(%v) detected. documents: %+v", col.Name(), time.Since(startTime), documents))
	}
	return insertOneResult, err
}

func (col *Collection[T]) updateOne(logger Logger, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	updateResult, err := col.Collection.UpdateOne(ctx, filter, update, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s UpdateOne slow query(%v) detected. filter: %+v, update: %+v", col.Name(), time.Since(startTime), filter, update))
	}
	return updateResult, err
}

func (col *Collection[T]) updateMany(logger Logger, ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	updateResult, err := col.Collection.UpdateMany(ctx, filter, update, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s updateMany slow query(%v) detected. filter: %+v, update: %+v", col.Name(), time.Since(startTime), filter, update))
	}
	return updateResult, err
}

func (col *Collection[T]) replaceOne(logger Logger, ctx context.Context, filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	result, err := col.Collection.ReplaceOne(ctx, filter, document, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s replaceOne slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return result, err
}

func (col *Collection[T]) deleteOne(logger Logger, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	startTime := time.Now()
	deleteResult, err := col.Collection.DeleteOne(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfOne() {
		logger.SlowQuery(fmt.Sprintf("%s deleteOne slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return deleteResult, err
}

func (col *Collection[T]) deleteMany(logger Logger, ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	startTime := time.Now()
	deleteResult, err := col.Collection.DeleteMany(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s deleteMany slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return deleteResult, err
}

func (col *Collection[T]) countDocuments(logger Logger, ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	startTime := time.Now()
	count, err := col.Collection.CountDocuments(ctx, filter, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s countDocuments slow query(%v) detected. filter: %+v", col.Name(), time.Since(startTime), filter))
	}
	return count, err
}

func (col *Collection[T]) estimatedDocumentCount(logger Logger, ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	startTime := time.Now()
	count, err := col.Collection.EstimatedDocumentCount(ctx, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfMany() {
		logger.SlowQuery(fmt.Sprintf("%s estimatedDocumentCount slow query(%v) detected.", col.Name(), time.Since(startTime)))
	}
	return count, err
}

func (col *Collection[T]) bulkWrite(logger Logger, ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	startTime := time.Now()
	bulkWriteResult, err := col.Collection.BulkWrite(ctx, models, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfBulk() {
		logger.SlowQuery(fmt.Sprintf("%s bulkWrite slow query(%v) detected. models: %+v", col.Name(), time.Since(startTime), models))
	}
	return bulkWriteResult, err
}

func (col *Collection[T]) aggregate(logger Logger, ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	startTime := time.Now()
	cursor, err := col.Collection.Aggregate(ctx, pipeline, opts...)
	if time.Since(startTime) >= logger.GetSlowQueryDurationOfAggregation() {
		logger.SlowQuery(fmt.Sprintf("%s aggregate slow query(%v) detected. pipeline: %+v", col.Name(), time.Since(startTime), pipeline))
	}
	return cursor, err
}
