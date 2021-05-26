package storytellermongo

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase struct {
	storiesCollection *mongo.Collection
	database          *mongo.Database
}

func (db *MongoDatabase) InitDB() error {
	// Database Config
	clientOptions := options.Client().ApplyURI(os.Getenv("CONNECTION_STRING"))
	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	//To close the connection at the end
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	} else {
		log.Println("Connected!")
	}
	db.database = client.Database(os.Getenv("POSTS_DB_NAME"))
	db.setCollections()
	return nil
}

//setCollections sets the db and correct collection
func (db *MongoDatabase) setCollections() {
	db.storiesCollection = db.database.Collection(os.Getenv("STORY_COLLECTION_NAME"))
}
