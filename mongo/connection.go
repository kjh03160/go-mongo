package mongo

import (
	"context"
	"fmt"
	"time"

	"mongo-orm/errorType"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_timeout = 10 * time.Second

type Client struct {
	authSource string
	*mongo.Client
}

var MongoClient *Client

func Connect(authSource string, timeout time.Duration, clientOpt *options.ClientOptions) *Client {
	c, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		panic(err.Error())
	}
	MongoClient = &Client{authSource, c}
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err.Error())
	}
	DB_timeout = timeout
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
