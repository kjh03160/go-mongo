package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"mongo-orm/data"
	"mongo-orm/errorType"
	"mongo-orm/mongo"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type MyLogger struct {
	*logrus.Logger
}

func (l *MyLogger) GetQueryTimeoutDuration() time.Duration {
	return 2 * time.Second
}

func (l *MyLogger) SlowQuery(msg string) {
	l.Error(msg)
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

func main() {
	c := mongo.Connect("sample_analytics", 10*time.Second, getMongoConfig())
	m := mongo.MakeCollection[data.Account](mongo.MongoClient, "sample_analytics", "accounts")
	defer c.Disconnect()

	logger := MyLogger{logrus.New()}
	r := gin.Default()
	r.GET("/find-one", func(context *gin.Context) {
		query := context.Query("account_id")
		accountId, _ := strconv.Atoi(query)
		var t data.Account
		if err := m.FindOne(&logger, &t, bson.M{"account_id": accountId}); err != nil {
			if errorType.IsNotFoundErr(err) {
				context.JSON(http.StatusNotFound, err.Error())
				return
			}
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, t)
	})

	r.GET("/find-all", func(context *gin.Context) {
		all, err := m.FindAll(&logger, bson.M{})
		if err != nil {
			if errorType.IsDecodeError(err) {
				fmt.Println("decode err")
				context.JSON(http.StatusInternalServerError, err.Error())
				return
			} else {
				fmt.Println("internal err")
				context.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		context.JSON(http.StatusOK, all)
	})

	r.GET("/insert-many", func(context *gin.Context) {
		var result []data.Account
		accounts, _ := m.FindAll(&logger, bson.M{})
		accounts = result[:2]
		var slice []interface{}
		for _, v := range accounts {
			slice = append(slice, v)
		}
		many, err := m.InsertMany(&logger, slice)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(many)
		context.JSON(http.StatusOK, result)
	})

	r.Run()
}
