package mongo

import (
	"os"
	"time"

	"github.com/kjh03160/go-mongo/errorType"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Account struct {
	AccountId int      `json:"account_id" bson:"account_id"`
	Limit     int      `json:"limit" bson:"limit"`
	Products  []string `json:"products" bson:"products"`
}

func getMongoConfig() *options.ClientOptions {
	mongoURI := os.Getenv("local")

	clientOptions := options.Client()
	clientOptions.ApplyURI(mongoURI).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMaxConnIdleTime(10 * time.Minute)
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.W(1)))
	clientOptions.SetReadConcern(readconcern.Local())
	clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	return clientOptions
}

func example_finoOne() {
	client := Connect(getMongoConfig())
	collection := NewCollection[Account](MongoClient, "sample_analytics", "accounts")
	defer client.Disconnect()

	logger := MyLogger{logrus.New()}
	accountId := 1

	var t Account
	if err := collection.FindOne(&logger, &t, bson.M{"account_id": accountId}); err != nil {
		if errorType.IsNotFoundErr(err) {
			logger.Warn(err.Error())
			return
		}
		logger.Error(err.Error())
		return
	}
	logger.Info(t)
}

func example_find_all() {
	client := Connect(getMongoConfig())
	collection := NewCollection[Account](MongoClient, "sample_analytics", "accounts")
	defer client.Disconnect()

	logger := MyLogger{logrus.New()}

	all, err := collection.FindAll(&logger, bson.M{})
	if err != nil {
		if errorType.IsDecodeError(err) {
			logger.Warn("decode err", err.Error())
			return
		} else {
			logger.Error(err.Error())
			return
		}
	}
	logger.Info(all)
}

func example_insert_many() {
	client := Connect(getMongoConfig())
	collection := NewCollection[Account](MongoClient, "sample_analytics", "accounts")
	defer client.Disconnect()

	logger := MyLogger{logrus.New()}

	var result []Account
	accounts, _ := collection.FindAll(&logger, bson.M{})

	accounts = result[:2]

	var slice []interface{}
	for _, v := range accounts {
		slice = append(slice, v)
	}

	insertedIds, err := collection.InsertMany(&logger, slice)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info(insertedIds)
}
