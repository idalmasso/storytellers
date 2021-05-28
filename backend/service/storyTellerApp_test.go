package service

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/idalmasso/storytellers/backend/common"
	"github.com/idalmasso/storytellers/backend/database/memorydb"
)

type testMemDb struct {
	stories      []common.Story
	testingValue int
}

func (db *testMemDb) InitDB(ctx context.Context) error {
	if db.testingValue == 0 {
		db.stories = make([]common.Story, 0)
	} else if db.testingValue == 1 {
		return fmt.Errorf("Error Init DB")
	}
	return nil
}
func (db *testMemDb) AddStory(ctx context.Context, story common.Story) (common.Story, error) {
	if db.testingValue != 2 {
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
		db.stories = append(db.stories, story)
		return story, nil
	}
	return common.Story{}, fmt.Errorf("%d", db.testingValue)
}
func (db *testMemDb) DeleteStory(ctx context.Context, id string) error {
	if db.testingValue == 3 {
		return fmt.Errorf("database error")
	}
	for i, v := range db.stories {
		if v.ID == id {
			db.stories[i] = db.stories[len(db.stories)-1]
			db.stories = db.stories[:len(db.stories)-1]
			return nil
		}
	}
	return fmt.Errorf("Not found")
}
func (db *testMemDb) UpdateStory(ctx context.Context, id string, story common.Story) (common.Story, error) {
	if db.testingValue == 4 {
		return common.Story{}, fmt.Errorf("database error")
	}
	for i, v := range db.stories {
		if v.ID == id {
			db.stories[i].Text = story.Text
			db.stories[i].Title = story.Title
			db.stories[i].User = story.User
			db.stories[i].ImagesPath = story.ImagesPath
			return db.stories[i], nil
		}
	}
	return common.Story{}, fmt.Errorf("Not found")
}
func (db *testMemDb) FindStory(ctx context.Context, id string) (common.Story, error) {
	if db.testingValue == 5 {
		return common.Story{}, fmt.Errorf("database error")
	}
	for i, v := range db.stories {
		if v.ID == id {
			return db.stories[i], nil
		}
	}
	return common.Story{}, fmt.Errorf("Not found")
}
func (db *testMemDb) FindAllStories(ctx context.Context) ([]common.Story, error) {
	if db.testingValue == 6 {
		return nil, fmt.Errorf("database error")
	}
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
			t.Error("No error expected")
		}
	})
	t.Run("App create with errors", func(t *testing.T) {
		db := testMemDb{testingValue: 1}
		app, err := CreateApp(&db)
		if err == nil || app.db != nil {
			t.Error("Error expected")
		}
	})
}

func TestCreateStory(t *testing.T) {
	t.Run("Create story ok", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)

		story, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("No error expected")
		}
		if story.Title != "title" || story.User != "user" {
			t.Error("Returned story wrong")
		}
		if story.ID == "" {
			t.Error("returned story with no id")
		}
		for _, s := range db.stories {
			if s.Title == "title" && s.User == "user" {
				return
			}
		}
		t.Error("cannot find the story in db")

	})
	t.Run("Create story ko with error db", func(t *testing.T) {
		db := testMemDb{testingValue: 2}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "user", "title")
		if err == nil {
			t.Error("error expected")
		}
		if story.Title != "" || story.User != "" {
			t.Error("returned story wrong, expected empty")
		}
		if story.ID != "" {
			t.Error("returned story with id")
		}

		if len(db.stories) != 0 {
			t.Error("There is a story in db, should not")
		}
	})
	t.Run("Create story with empty title", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "user", "")
		if err == nil {
			t.Error("error expected")
		}
		if story.Title != "" || story.User != "" {
			t.Error("returned story wrong, expected empty")
		}
		if story.ID != "" {
			t.Error("returned story with id")
		}

		if len(db.stories) != 0 {
			t.Error("There is a story in db, should not")
		}

	})
	t.Run("Create story with empty user", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "", "title")
		if err == nil {
			t.Error("error expected")
		}
		if story.ID != "" {
			t.Error("returned story with id")
		}
		if story.Title != "" || story.User != "" {
			t.Error("returned story wrong, expected empty")
		}

		if len(db.stories) != 0 {
			t.Error("There is a story in db, should not")
		}

	})
}

func TestDeleteStory(t *testing.T) {
	t.Run("Delete story ok", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		err = app.DeleteStory(context.Background(), story.ID)
		if err != nil {
			t.Error("error returned from delete story")
		}
		for _, s := range db.stories {
			if s.Title == "title" && s.User == "user" {
				t.Error("record present in database, not expected")
			}
		}
	})
	t.Run("Delete story ko with error db", func(t *testing.T) {
		db := testMemDb{testingValue: 3}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		err = app.DeleteStory(context.Background(), s.ID)
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		for _, s := range db.stories {
			if s.Title == "title" && s.User == "user" {
				return
			}
		}
		t.Error("record not present in database, expected after error")
	})
	t.Run("Delete story with empty id", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		err = app.DeleteStory(context.Background(), "")
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		for _, s := range db.stories {
			if s.Title == "title" && s.User == "user" {
				return

			}
		}
		t.Error("record not present in database, expected after error")
	})
}

func TestFindStory(t *testing.T) {
	t.Run("Find story ok", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.FindStory(context.Background(), s.ID)
		if err != nil {
			t.Error("error returned from find story")
		}
		if story.User != "user" || story.Title != "title" {
			t.Error("story returned wrong")
		}

	})
	t.Run("Find story ko with error db", func(t *testing.T) {
		db := testMemDb{testingValue: 5}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.FindStory(context.Background(), s.ID)
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		if story.Title != "" || story.User != "" {
			t.Error("story returned wrong")
		}
	})
	t.Run("Find story with empty id", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		story, err := app.FindStory(context.Background(), "")
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		if story.Title != "" || story.User != "" {
			t.Error("story returned wrong")
		}
	})
}

func TestFindAllStories(t *testing.T) {
	t.Run("FindAllStories ok one story", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 1 {
			t.Error("stories returned wrong")
		}

	})
	t.Run("FindAllStories ko with error db", func(t *testing.T) {
		db := testMemDb{testingValue: 6}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		stories, err := app.FindAllStories(context.Background())
		if err == nil {
			t.Error("expected error returned from FindAllStories")
		}
		if stories != nil {
			t.Error("stories returned wrong")
		}
	})
	t.Run("FindAllStories ok two stories", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		_, err = app.CreateStory(context.Background(), "user2", "title2")
		if err != nil {
			t.Error("error returned from create story2")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 2 {
			t.Error("stories returned wrong")
		}

	})
}

func TestUpdateStory(t *testing.T) {
	t.Run("Update story ok", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.UpdateStory(context.Background(), s.ID, "user", "title", "my text")
		if err != nil {
			t.Error("error returned from UpdateStory")
		}
		if story.User != "user" || story.Title != "title" || story.Text != "my text" {
			t.Error("story returned wrong")
		}
		story, _ = app.FindStory(context.Background(), s.ID)
		if story.User != "user" || story.Title != "title" || story.Text != "my text" {
			t.Error("story in db wrong")
		}
	})
	t.Run("Update story ko with error db", func(t *testing.T) {
		db := testMemDb{testingValue: 4}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.UpdateStory(context.Background(), s.ID, "user", "title", "my text")
		if err == nil {
			t.Error("error expected from UpdateStory")
		}
		if story.User != "" || story.Title != "" || story.Text != "" {
			t.Error("story returned wrong")
		}
		story, _ = app.FindStory(context.Background(), s.ID)
		if story.User != "user" || story.Title != "title" || story.Text != "" {
			t.Error("story in db wrong")
		}
	})
	t.Run("Update story with empty id", func(t *testing.T) {
		db := testMemDb{testingValue: 0}
		app, _ := CreateApp(&db)

		story, err := app.UpdateStory(context.Background(), "", "user", "", "my text")
		if err == nil {
			t.Error("error expected from UpdateStory")
		}
		if story.User != "" || story.Title != "" || story.Text != "" {
			t.Error("story returned wrong")
		}
	})
}

///Here are tests for real DB, integration tests!
func TestAppCreate_Integration(t *testing.T) {
	t.Run("App create with no errors memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, err := CreateApp(&db)
		if err != nil || app.db == nil {
			t.Error("No error expected")
		}
	})

}

func TestCreateStory_Integration(t *testing.T) {
	t.Run("Create story ok memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)

		story, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("No error expected")
		}
		if story.Title != "title" || story.User != "user" {
			t.Error("Returned story wrong")
		}
		if story.ID == "" {
			t.Error("returned story with no id")
		}
		stories, err := db.FindAllStories(context.Background())
		if err != nil {
			t.Error("returned error in findallstories " + err.Error())
		}
		for _, s := range stories {
			if s.Title == "title" && s.User == "user" {
				return
			}
		}
		t.Error("cannot find the story in db")

	})
	t.Run("Create story with empty title memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "user", "")
		if err == nil {
			t.Error("error expected")
		}
		if story.Title != "" || story.User != "" {
			t.Error("returned story wrong, expected empty")
		}
		if story.ID != "" {
			t.Error("returned story with id")
		}
		stories, err := db.FindAllStories(context.Background())
		if err != nil {
			t.Error("returned error in findallstories " + err.Error())
		}
		if len(stories) != 0 {
			t.Error("There is a story in db, should not")
		}

	})
	t.Run("Create story with empty user memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "", "title")
		if err == nil {
			t.Error("error expected")
		}
		if story.ID != "" {
			t.Error("returned story with id")
		}
		if story.Title != "" || story.User != "" {
			t.Error("returned story wrong, expected empty")
		}
		stories, err := db.FindAllStories(context.Background())
		if err != nil {
			t.Error("returned error in findallstories " + err.Error())
		}
		if len(stories) != 0 {
			t.Error("There is a story in db, should not")
		}

	})
}

func TestDeleteStory_Integration(t *testing.T) {
	t.Run("Delete story ok memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		story, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		err = app.DeleteStory(context.Background(), story.ID)
		if err != nil {
			t.Error("error returned from delete story")
		}
		stories, err := db.FindAllStories(context.Background())
		if err != nil {
			t.Error("returned error in findallstories " + err.Error())
		}
		for _, s := range stories {
			if s.Title == "title" && s.User == "user" {
				t.Error("record present in database, not expected")
			}
		}
	})

	t.Run("Delete story with empty id memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		err = app.DeleteStory(context.Background(), "")
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		stories, err := db.FindAllStories(context.Background())
		if err != nil {
			t.Error("returned error in findallstories " + err.Error())
		}
		for _, s := range stories {
			if s.Title == "title" && s.User == "user" {
				return
			}
		}
		t.Error("record not present in database, expected after error")
	})
}

func TestFindStory_Integration(t *testing.T) {
	t.Run("Find story ok memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.FindStory(context.Background(), s.ID)
		if err != nil {
			t.Error("error returned from find story")
		}
		if story.User != "user" || story.Title != "title" {
			t.Error("story returned wrong")
		}

	})
	t.Run("Find story with empty id memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		story, err := app.FindStory(context.Background(), "")
		if err == nil {
			t.Error("expected error returned from delete story")
		}
		if story.Title != "" || story.User != "" {
			t.Error("story returned wrong")
		}
	})
}

func TestFindAllStories_Integration(t *testing.T) {
	t.Run("FindAllStories ok one story memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 1 {
			t.Error("stories returned wrong")
		}

	})

	t.Run("FindAllStories ok two stories memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		_, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
		}
		_, err = app.CreateStory(context.Background(), "user2", "title2")
		if err != nil {
			t.Error("error returned from create story2")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 2 {
			t.Error("stories returned wrong")
		}

	})
}

func TestUpdateStory_Integration(t *testing.T) {
	t.Run("Update story ok memoryDB", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)
		s, err := app.CreateStory(context.Background(), "user", "title")
		if err != nil {
			t.Error("error returned from create story")
			return
		}
		story, err := app.UpdateStory(context.Background(), s.ID, "user", "title", "my text")
		if err != nil {
			t.Error("error returned from UpdateStory")
		}
		if story.User != "user" || story.Title != "title" || story.Text != "my text" {
			t.Error("story returned wrong")
		}
		story, _ = app.FindStory(context.Background(), s.ID)
		if story.User != "user" || story.Title != "title" || story.Text != "my text" {
			t.Error("story in db wrong")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 1 {
			t.Error("stories returned wrong")
		}
	})
	t.Run("Update story with memoryDB no story created", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)

		story, err := app.UpdateStory(context.Background(), "1", "user", "", "my text")
		if err == nil {
			t.Error("error expected from UpdateStory")
		}
		if story.User != "" || story.Title != "" || story.Text != "" {
			t.Error("story returned wrong")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 0 {
			t.Error("stories returned wrong")
		}
	})
	t.Run("Update story with memoryDB empty id", func(t *testing.T) {
		db := memorydb.MemoryDatabase{}
		app, _ := CreateApp(&db)

		story, err := app.UpdateStory(context.Background(), "", "user", "", "my text")
		if err == nil {
			t.Error("error expected from UpdateStory")
		}
		if story.User != "" || story.Title != "" || story.Text != "" {
			t.Error("story returned wrong")
		}
		stories, err := app.FindAllStories(context.Background())
		if err != nil {
			t.Error("got error returned from findAll, not expected")
		}
		if stories == nil || len(stories) != 0 {
			t.Error("stories returned wrong")
		}
	})
}
