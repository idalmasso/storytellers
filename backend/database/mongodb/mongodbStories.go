package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/idalmasso/storytellers/backend/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *MongoDatabase) AddStory(ctx context.Context, story common.Story) (common.Story, error) {
	result, err := db.storiesCollection.InsertOne(ctx, story)
	if err != nil {
		return common.Story{}, err
	}
	if val, ok := result.InsertedID.(primitive.ObjectID); ok {
		story.ID = val.Hex()
		return story, nil
	}
	return story, fmt.Errorf("Cannot get id from results")
}

func (db *MongoDatabase) DeleteStory(ctx context.Context, id string) error {
	//objID, err := primitive.ObjectIDFromHex(postID)
	realId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("cannot covert id " + err.Error())
	}
	result, err := db.storiesCollection.DeleteOne(ctx, bson.M{"_id": realId})
	if err != nil {
		return fmt.Errorf("cannot delete request " + err.Error())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("not found story")
	}
	return nil
}
func (db *MongoDatabase) UpdateStory(ctx context.Context, id string, storyUpdated common.Story) (common.Story, error) {
	realId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return common.Story{}, fmt.Errorf("cannot covert id " + err.Error())
	}
	result := db.storiesCollection.FindOneAndReplace(ctx, bson.M{"_id": realId}, storyUpdated)
	if result.Err() != nil {
		return common.Story{}, fmt.Errorf("Error, cannot update request " + result.Err().Error())
	}

	return common.Story{}, nil
}
func (db *MongoDatabase) FindStory(ctx context.Context, id string) (common.Story, error) {
	realId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return common.Story{}, fmt.Errorf("cannot covert id " + err.Error())
	}
	result := db.storiesCollection.FindOne(ctx, bson.M{"_id": realId})
	if result.Err() != nil {
		return common.Story{}, fmt.Errorf("Error, cannot find request " + result.Err().Error())
	}
	var story common.Story
	err = result.Decode(&story)
	if err != nil {
		return common.Story{}, fmt.Errorf("Error, cannot decode from database " + err.Error())
	}
	return story, nil
}
func (db *MongoDatabase) FindAllStories(ctx context.Context) ([]common.Story, error) {
	var stories []common.Story
	cursor, err := db.storiesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error while getting cursor, Reason: %v\n", err)
		return nil, err
	}
	stories = make([]common.Story, 0)
	for cursor.Next(ctx) {
		var story common.Story
		cursor.Decode(&story)
		stories = append(stories, story)
	}
	return stories, nil
}
