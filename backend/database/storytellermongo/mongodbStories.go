package storytellermongo

import (
	"context"
	"fmt"
	"log"

	"github.com/idalmasso/storytellers/backend/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *MongoDatabase) AddStory(story common.Story) (common.Story, error) {
	result, err := db.storiesCollection.InsertOne(context.TODO(), story)
	if err != nil {
		return story, err
	}
	if _, ok := result.InsertedID.(primitive.ObjectID); ok {
		//story.ID = oid
		return story, nil
	}
	return story, fmt.Errorf("Cannot get id from results")
}
func (db *MongoDatabase) DeleteStory(user, title string) error {
	//objID, err := primitive.ObjectIDFromHex(postID)
	result, err := db.storiesCollection.DeleteOne(context.TODO(), bson.M{"user": user, "title": title})
	if err != nil {
		return fmt.Errorf("Error, cannot delete request " + err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Not found story")
	}
	return nil
}
func (db *MongoDatabase) UpdateStory(story common.Story) (common.Story, error) {
	result := db.storiesCollection.FindOneAndReplace(context.TODO(), bson.M{"user": story.User, "title": story.Title}, story)
	if result.Err() != nil {
		return story, fmt.Errorf("Error, cannot update request " + result.Err().Error())
	}

	return story, nil
}
func (db *MongoDatabase) FindStory(user, title string) (common.Story, error) {
	result := db.storiesCollection.FindOne(context.TODO(), bson.M{"user": user, "title": title})
	if result.Err() != nil {
		return common.Story{}, fmt.Errorf("Error, cannot find request " + result.Err().Error())
	}
	var story common.Story
	err := result.Decode(&story)
	if err != nil {
		return common.Story{}, fmt.Errorf("Error, cannot decode from database " + err.Error())
	}
	return story, nil
}
func (db *MongoDatabase) FindAllStories() ([]common.Story, error) {
	var stories []common.Story
	cursor, err := db.storiesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error while getting cursor, Reason: %v\n", err)
		return nil, err
	}
	stories = make([]common.Story, 0)
	for cursor.Next(context.TODO()) {
		var story common.Story
		cursor.Decode(&story)
		stories = append(stories, story)
	}
	return stories, nil
}
