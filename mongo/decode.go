package mongo

import (
	"context"
	"fmt"
	"reflect"

	"mongo-orm/errorType"

	"go.mongodb.org/mongo-driver/mongo"
)

func EvaluateAndDecodeSingleResult(result *mongo.SingleResult, v interface{}) error {
	if result == nil {
		return errorType.SingleResultErr
	}
	if err := result.Decode(v); err != nil {
		return err
	}
	return nil
}

func DecodeCursor(cursor *mongo.Cursor, t reflect.Type) interface{} {
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

func DecodeCursorAll(cursor *mongo.Cursor, slice interface{}) error {
	if err := cursor.All(context.Background(), slice); err != nil {
		return err
	}
	return nil
}
