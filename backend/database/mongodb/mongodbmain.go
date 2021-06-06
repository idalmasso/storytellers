package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase struct {
	storiesCollection *mongo.Collection
	database          *mongo.Database
}

func (db *MongoDatabase) InitDB(ctx context.Context) error {
	// Database Config
	clientOptions := options.Client().ApplyURI(os.Getenv("CONNECTION_STRING"))
	client, err := mongo.NewClient(clientOptions)
	if err!=nil{
		return err
	}
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	} else {
		log.Println("Connected!")
	}
	db.database = client.Database(os.Getenv("STORIES_DB_NAME"))
	db.setCollections()
	return nil
}

//setCollections sets the db and correct collection
func (db *MongoDatabase) setCollections() {
	db.storiesCollection = db.database.Collection(os.Getenv("STORY_COLLECTION_NAME"))
}
