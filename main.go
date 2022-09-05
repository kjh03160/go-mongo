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
	c := mongo.Connect(mongoURI, "sample_analytics", mongo.SECONDARY_PREFERRED, 1)
	m := mongo.MakeCollection[data.Account](mongo.MongoClient, "sample_analytics", "accounts")
	defer c.Disconnect()

	r := gin.Default()
	r.GET("/find-one", func(context *gin.Context) {
		query := context.Query("account_id")
		accountId, _ := strconv.Atoi(query)
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
		all, err := m.FindAll(bson.M{})
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
		accounts, _ := m.FindAll(bson.M{})
		accounts = result[:2]
		var slice []interface{}
		for _, v := range accounts {
			slice = append(slice, v)
		}
		many, err := m.InsertMany(slice)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(many)
		context.JSON(http.StatusOK, result)
	})

	r.Run()
}
