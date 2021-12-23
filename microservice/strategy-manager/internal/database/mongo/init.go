package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DbClient   *mongo.Client
	MongoCtx   *context.Context
	StrategyDB *mongo.Collection
)

func MongoConnect() {
	mongoCtx := context.Background()
	MongoCtx = &mongoCtx
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		log.Println("Connected to Mongodb")
	}
	DbClient = db
	StrategyDB = DbClient.Database("mydb").Collection("strategy")
}
