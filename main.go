package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	URL        string = "127.0.0.1"
	USERNAME   string = "root"
	PASSWORD   string = "liurui97128224"
	PORT       uint   = 27017
	DATABASE   string = "test"
	COLLECTION string = "cron_log"
)

func main() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authMechanism=SCRAM-SHA-1", USERNAME, PASSWORD, URL, PORT, DATABASE))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	database := client.Database(DATABASE)
	collection := database.Collection(COLLECTION)
	fmt.Println(collection)
}
