package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"mongo-orm/data"
	"mongo-orm/errorType"
	"mongo-orm/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	mongoURI := os.Getenv("local")
	mongo.Connect(mongoURI, "sample_analytics", mongo.SECONDARY_PREFERRED)
	collection := mongo.GetCollection("accounts")
	defer mongo.Disconnect()
	m := &mongo.Collection{Collection: collection}

	r := SetRouter()
	r.GET("/find-one", func(context *gin.Context) {
		query := context.Query("account_id")
		accountId, _ := strconv.Atoi(query)
		var t data.Account
		if err := m.FindOne(&t, bson.M{"account_id": accountId}); err != nil {
			// 1. error catching by type switch
			if errorType.IsNotFoundErr(err) {
				context.JSON(http.StatusNotFound, err.Error())
				return
			}
			// 2. error catching using reflect
			if errorType.IsErrorTypeOf(err, errorType.GetNotFoundErrorType()) {
				context.JSON(http.StatusNotFound, err.Error())
				return
			}
			context.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		context.JSON(http.StatusOK, t)
	})

	r.GET("/find-all", func(context *gin.Context) {
		all, _ := m.FindAll(data.Account{}, bson.M{})
		context.JSON(http.StatusOK, all.([]data.Account))
	})

	r.GET("/find-all/v2", func(context *gin.Context) {
		var result []data.Account
		err := m.FindAllV2(&result, bson.M{})
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
		context.JSON(http.StatusOK, result)
	})

	r.Run()
}
