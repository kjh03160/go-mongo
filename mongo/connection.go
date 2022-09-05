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
	authSource string
	*mongo.Client
}

var MongoClient *Client

func Connect(uri string, authSource, readPreference string, maxPoolSize uint64) *Client {
	clientOptions := options.Client()
	clientOptions.ApplyURI(uri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))

	clientOptions.SetMaxPoolSize(maxPoolSize)
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
		panic(err.Error())
	}
	MongoClient = &Client{authSource, c}
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err.Error())
	}
	return MongoClient
}

func (client *Client) Disconnect() {
	err := client.Client.Disconnect(context.TODO())
	if err != nil {
		panic(nil)
	}
	fmt.Printf("Connections to MongoDB %s closed\n", client.authSource)
}

func (client *Client) GetDatabase(dbName string) *mongo.Database {
	return client.Database(dbName)
}

func (client *Client) GetCollection(databaseName, collection string) *mongo.Collection {
	return client.Database(databaseName).Collection(collection)
}

func (client *Client) Transaction(sessionOpt *options.SessionOptions, trxOpt *options.TransactionOptions, function func(sessCtx mongo.SessionContext) (interface{}, error)) error {
	if sessionOpt == nil {
		sessionOpt = &options.SessionOptions{}
	}
	if trxOpt == nil {
		trxOpt = &options.TransactionOptions{}
	}

	session, err := client.Client.StartSession(sessionOpt)
	if err != nil {
		return errorType.MongoClientError(err)
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), function, trxOpt)
	if err != nil {
		return err
	}
	return nil
}
