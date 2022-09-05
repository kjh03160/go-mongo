package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[T any] struct {
	*mongo.Collection
}

func MakeCollection[T any](mongoManager *Client, databaseName, collectionName string) *Collection[T] {
	collection := mongoManager.GetCollection(databaseName, collectionName)
	return &Collection[T]{Collection: collection}
}

func (col *Collection[T]) findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOne(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return singleResult
}

func (col *Collection[T]) findAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	startTime := time.Now()
	cursor, err := col.Collection.Find(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return cursor, err
}

func (col *Collection[T]) findOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	fmt.Println(time.Since(startTime))
	return singleResult
}

func (col *Collection[T]) findOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndReplace(ctx, filter, replacement, opts...)
	fmt.Println(time.Since(startTime))
	return singleResult
}

func (col *Collection[T]) findOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	startTime := time.Now()
	singleResult := col.Collection.FindOneAndDelete(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return singleResult
}

func (col *Collection[T]) insertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	startTime := time.Now()
	insertOneResult, err := col.Collection.InsertOne(ctx, document, opts...)
	fmt.Println(time.Since(startTime))
	return insertOneResult, err
}

func (col *Collection[T]) insertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	startTime := time.Now()
	insertOneResult, err := col.Collection.InsertMany(ctx, documents, opts...)
	fmt.Println(time.Since(startTime))
	return insertOneResult, err
}

func (col *Collection[T]) updateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	updateResult, err := col.Collection.UpdateOne(ctx, filter, update, opts...)
	fmt.Println(time.Since(startTime))
	return updateResult, err
}

func (col *Collection[T]) updateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	updateResult, err := col.Collection.UpdateOne(ctx, filter, update, opts...)
	fmt.Println(time.Since(startTime))
	return updateResult, err
}

func (col *Collection[T]) replaceOne(ctx context.Context, filter interface{}, document interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	result, err := col.Collection.ReplaceOne(ctx, filter, document, opts...)
	fmt.Println(time.Since(startTime))
	return result, err
}

func (col *Collection[T]) deleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	startTime := time.Now()
	deleteResult, err := col.Collection.DeleteOne(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return deleteResult, err
}

func (col *Collection[T]) deleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	startTime := time.Now()
	deleteResult, err := col.Collection.DeleteMany(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return deleteResult, err
}

func (col *Collection[T]) countDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	startTime := time.Now()
	count, err := col.Collection.CountDocuments(ctx, filter, opts...)
	fmt.Println(time.Since(startTime))
	return count, err
}

func (col *Collection[T]) estimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	startTime := time.Now()
	count, err := col.Collection.EstimatedDocumentCount(ctx, opts...)
	fmt.Println(time.Since(startTime))
	return count, err
}

func (col *Collection[T]) bulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	startTime := time.Now()
	bulkWriteResult, err := col.Collection.BulkWrite(ctx, models, opts...)
	fmt.Println(time.Since(startTime))
	return bulkWriteResult, err
}

func (col *Collection[T]) aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	startTime := time.Now()
	cursor, err := col.Collection.Aggregate(ctx, pipeline, opts...)
	fmt.Println(time.Since(startTime))
	return cursor, err
}
