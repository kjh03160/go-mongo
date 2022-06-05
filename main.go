package main

import (
	"mongo-orm/connection"
	"os"
)

func main() {
	db := connection.MongoDBManager{}
	mongoURI := os.Getenv("local")
	db.Connect(mongoURI, "admin", connection.SECONDARY_PREFERRED)

	r := SetRouter()
	r.Run()
}
