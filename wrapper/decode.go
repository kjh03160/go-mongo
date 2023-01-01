package wrapper

import (
	"context"

	"github.com/kjh03160/go-mongo/errorType"
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

func DecodeCursor[T any](cursor *mongo.Cursor) ([]T, error) {
	defer cursor.Close(context.Background())
	slice := []T{}
	for cursor.Next(context.Background()) {
		var doc T
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		slice = append(slice, doc)
	}
	return slice, nil
}
