package service

import "github.com/idalmasso/storytellers/backend/common"

type storyTellerDB interface {
	InitDB() error
	AddStory(common.Story) (common.Story, error)
	DeleteStory(user, title string) error
	UpdateStory(common.Story) (common.Story, error)
	FindStory(user, title string) (common.Story, error)
	FindAllStories() ([]common.Story, error)
}

type storyTellerApp struct {
	db storyTellerDB
}

func CreateApp(db storyTellerDB) (storyTellerApp, error) {
	app := storyTellerApp{db: db}
	err := app.db.InitDB()
	if err != nil {
		return storyTellerApp{}, err
	}
	return app, nil
}

func (app storyTellerApp) CreateStory(user, title string) (common.Story, error) {
	story := common.Story{User: user, Title: title, Text: "", ImagesPath: ""}
	story, err := app.db.AddStory(story)
	if err != nil {
		return common.Story{}, err
	}
	return story, nil
}

func (app storyTellerApp) DeleteStory(user, title string) error {
	return app.db.DeleteStory(user, title)
}

func (app storyTellerApp) FindStory(user, title string) (common.Story, error) {
	return app.db.FindStory(user, title)
}
func (app storyTellerApp) FindAllStories() ([]common.Story, error) {
	return app.db.FindAllStories()
}
