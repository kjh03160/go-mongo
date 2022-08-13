package mongo

import (
	"context"
	"fmt"
	"time"

	"mongo-orm/errorType"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const PRIMARY = "PRIMARY"
const SECONDARY_PREFERRED = "SECONDARY_PREFERRED"
const PRIMARY_PREFERRED = "PRIMARY_PREFERRED"

const DB_TIMEOUT = 10 * time.Second
const DB_TIMEOUT_MANY = 20 * time.Second

type Client struct {
	*mongo.Client
}

var MongoClient *Client
var db *mongo.Database

func Connect(uri string, authSource, readPreference string) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(uri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	clientOptions.SetMaxConnIdleTime(10 * time.Minute)
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.W(1)))
	clientOptions.SetReadConcern(readconcern.Local())

	if readPreference == PRIMARY {
		clientOptions.SetReadPreference(readpref.Primary())
	} else if readPreference == PRIMARY_PREFERRED {
		clientOptions.SetReadPreference(readpref.PrimaryPreferred())
	} else if readPreference == SECONDARY_PREFERRED {
		clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	} else {
		clientOptions.SetReadPreference(readpref.SecondaryPreferred())
	}

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		panic(nil)
	}
	MongoClient = &Client{c}
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		panic(nil)
	}

	db = MongoClient.Database(authSource)

	fmt.Printf("Connections to MongoDB %s %s\n", uri, authSource)
}

func Disconnect() {
	err := MongoClient.Disconnect(context.TODO())
	if err != nil {
		panic(nil)
	}
	fmt.Printf("Connections to MongoDB %s closed\n", db.Name())
}

func (manager *Client) GetDatabase() *mongo.Database {
	return db
}

func (manager *Client) GetCollection(databaseName, collection string) *mongo.Collection {
	return manager.Database(databaseName).Collection(collection)
}

func (manager *Client) Transaction(readPreference *readpref.ReadPref, function func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	session, err := manager.Client.StartSession()
	if err != nil {
		return errorType.MongoClientError(err)
	}
	defer session.EndSession(context.Background())

	opt := options.TransactionOptions{}
	opt.SetReadPreference(readPreference)
	_, err = session.WithTransaction(context.Background(), function, &opt)
	if err != nil {
		return err
	}
	return nil
}
