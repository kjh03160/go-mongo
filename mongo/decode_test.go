package mongo

import (
	"context"
	"testing"

	"mongo-orm/data"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_EvaluateAndDecodeSingleResult(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(t *mtest.T) {
		col := t.Coll

		expected := data.Account{
			AccountId: 1,
		}
		t.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", primitive.NewObjectID()},
			{"account_id", expected.AccountId},
		}))

		singleResult := col.FindOne(context.Background(), bson.M{})

		var w data.Account
		err := EvaluateAndDecodeSingleResult(singleResult, &w)
		assert.NoError(t, err)
		assert.Equal(t, expected, w)
	})

	mt.Run("not found", func(t *mtest.T) {
		singleResult := mongo.NewSingleResultFromDocument(data.Account{}, mongo.ErrNoDocuments, nil)
		var w data.Account
		err := EvaluateAndDecodeSingleResult(singleResult, &w)
		assert.Error(t, err)
	})

	mt.Run("err", func(t *mtest.T) {
		singleResult := mongo.NewSingleResultFromDocument(nil, errors.New("test"), nil)
		var w data.Account
		err := EvaluateAndDecodeSingleResult(singleResult, &w)
		assert.Error(t, err)
	})
}

func Test_DecodeCursor(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(t *mtest.T) {
		col := t.Coll

		find := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"account_id", 1},
			})
		getMore := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch,
			bson.D{
				{"_id", primitive.NewObjectID()},
				{"account_id", 2},
			})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		t.AddMockResponses(find, getMore, killCursors)
		cursor, err := col.Find(context.Background(), bson.M{})

		resultSlice, err := DecodeCursor[data.Account](cursor)
		assert.NoError(t, err)
		assert.NotEmpty(t, resultSlice)
	})

	mt.Run("success - find result empty", func(t *mtest.T) {
		col := t.Coll

		find := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch)
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		t.AddMockResponses(find, killCursors)

		cursor, err := col.Find(context.Background(), bson.M{})

		resultSlice, err := DecodeCursor[data.Account](cursor)
		assert.NoError(t, err)
		assert.Empty(t, resultSlice)
		assert.NotNil(t, resultSlice)
	})
}
