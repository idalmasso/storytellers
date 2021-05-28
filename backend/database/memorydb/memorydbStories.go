package memorydb

import (
	"context"
	"fmt"
	"strconv"

	"github.com/idalmasso/storytellers/backend/common"
)

func (db *MemoryDatabase) AddStory(ctx context.Context, story common.Story) (common.Story, error) {
	var maxValue int64
	for _, s := range db.stories {
		val, err := strconv.ParseInt(s.ID, 10, 64)
		if err != nil {
			return common.Story{}, fmt.Errorf("cannot convert %v to int ", s.ID)
		}
		if val > maxValue {
			maxValue = val
		}
	}
	story.ID = strconv.FormatInt(maxValue+1, 10)
	db.stories[story.ID] = story
	return story, nil

}
func (db *MemoryDatabase) DeleteStory(ctx context.Context, id string) error {
	if _, ok := db.stories[id]; ok {
		delete(db.stories, id)
		return nil
	} else {
		return fmt.Errorf("Not found")
	}
}
func (db *MemoryDatabase) UpdateStory(ctx context.Context, id string, story common.Story) (common.Story, error) {
	if _, ok := db.stories[id]; ok {
		db.stories[id] = story
		return db.stories[id], nil
	} else {
		return common.Story{}, fmt.Errorf("Not found")
	}
}
func (db *MemoryDatabase) FindStory(ctx context.Context, id string) (common.Story, error) {
	if _, ok := db.stories[id]; ok {
		return db.stories[id], nil
	} else {
		return common.Story{}, fmt.Errorf("Not found")
	}
}
func (db *MemoryDatabase) FindAllStories(ctx context.Context) ([]common.Story, error) {
	stories := make([]common.Story, len(db.stories))
	i := 0
	for _, v := range db.stories {
		stories[i] = v
		i++
	}
	return stories, nil
}
