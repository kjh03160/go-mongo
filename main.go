package main

import (
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
	defer mongo.Disconnect()
	r := SetRouter()
	r.GET("/find-one", func(context *gin.Context) {
		query := context.Query("account_id")
		accountId, _ := strconv.Atoi(query)
		collection := mongo.GetCollection("accounts")
		m := &mongo.Collection{Collection: collection}
		var t data.Account
		if err := m.FindOne(&t, bson.M{"account_id": accountId}); err != nil {
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
		collection := mongo.GetCollection("accounts")
		m := &mongo.Collection{Collection: collection}
		all, _ := m.FindAll(data.Account{}, bson.M{})
		context.JSON(http.StatusOK, all.([]data.Account))
	})
	r.Run()
}
