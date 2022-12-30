package mongo

import (
	"context"
	"fmt"

	"github.com/kjh03160/go-mongo/errorType"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	*mongo.Client
}

var MongoClient *Client

func Connect(clientOpt *options.ClientOptions) *Client {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	MongoClient = &Client{Client: client}
	fmt.Println("Successfully connected and pinged.")
	return MongoClient
}

func (client *Client) Disconnect() {
	err := client.Client.Disconnect(context.TODO())
	if err != nil {
		panic(nil)
	}
	fmt.Println("Connections to MongoDB closed")
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
