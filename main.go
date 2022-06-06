package main

import (
	"mongo-orm/mongo"
	"os"
)

func main() {
	mongoURI := os.Getenv("local")
	mongo.Connect(mongoURI, "local", mongo.SECONDARY_PREFERRED)
	defer mongo.Disconnect()
	r := SetRouter()
	r.Run()
}
