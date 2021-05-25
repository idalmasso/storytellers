package service

import (
	"fmt"
	"testing"

	"github.com/idalmasso/storytellers/backend/common"
)

type testMemDb struct {
	stories      []common.Story
	testingValue int
}

func (db *testMemDb) InitDB() error {
	if db.testingValue == 0 {
		db.stories = make([]common.Story, 0)
	} else {
		return fmt.Errorf("Error Init DB")
	}
	return nil
}

func (db *testMemDb) AddStory(story common.Story) (common.Story, error) {
	db.stories = append(db.stories, story)
	if db.testingValue == 0 {
		return story, nil
	}
	return common.Story{}, fmt.Errorf("%d", db.testingValue)
}
func (db *testMemDb) DeleteStory(user, title string) error {
	for i, v := range db.stories {
		if v.User == user && v.Title == title {
			db.stories[i] = db.stories[len(db.stories)-1]
			db.stories = db.stories[:len(db.stories)-1]
			return nil
		}
	}
	return fmt.Errorf("Not found")
}
func (db *testMemDb) UpdateStory(story common.Story) (common.Story, error) {
	for i, v := range db.stories {
		if v.User == story.User && v.Title == story.Title {
			db.stories[i] = story
			return db.stories[i], nil
		}
	}
	return common.Story{}, fmt.Errorf("Not found")
}
func (db *testMemDb) FindStory(user, title string) (common.Story, error) {
	for i, v := range db.stories {
		if v.User == user && v.Title == title {
			return db.stories[i], nil
		}
	}
	return common.Story{}, fmt.Errorf("Not found")
}
func (db *testMemDb) FindAllStories() ([]common.Story, error) {
	if db.testingValue == 0 {
		return db.stories, nil
	}
	return nil, fmt.Errorf("%d", db.testingValue)
}
func TestAppCreate(t *testing.T) {

	t.Run("App create with no errors", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, err := CreateApp(&db)
		if err != nil || app.db == nil {
			t.Error("No error expected!")
		}
	})
	t.Run("App create with errors", func(t *testing.T) {
		db := testMemDb{testingValue: 1}
		app, err := CreateApp(&db)
		if err == nil || app.db != nil {
			t.Error("Error expected!")
		}
	})
}
