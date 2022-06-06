package mongo

import (
	"context"
	"fmt"
	"time"

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

var client *mongo.Client
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

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		panic(nil)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
		panic(nil)
	}

	db = client.Database(authSource)

	fmt.Printf("Connections to MongoDB %s %s\n", uri, authSource)
}

func Disconnect() {
	err := client.Disconnect(context.TODO())
	if err != nil {
		panic(nil)
	}
	fmt.Printf("Connections to MongoDB %s closed\n", db.Name())
}

func GetDatabase() *mongo.Database {
	return db
}

func GetCollection(collection string, opts ...*options.CollectionOptions) *mongo.Collection {
	return db.Collection(collection, opts...)
}
