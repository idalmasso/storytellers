package service

import (
	"context"
	"fmt"
	"time"

	"github.com/idalmasso/storytellers/backend/common"
)

type storyTellerDB interface {
	InitDB(context.Context) error
	AddStory(context.Context, common.Story) (common.Story, error)
	DeleteStory(ctx context.Context, id string) error
	UpdateStory(ctx context.Context, id string, updatedStory common.Story) (common.Story, error)
	FindStory(ctx context.Context, id string) (common.Story, error)
	FindAllStories(context.Context) ([]common.Story, error)
}

type storyTellerApp struct {
	db storyTellerDB
}

func CreateApp(db storyTellerDB) (storyTellerApp, error) {
	app := storyTellerApp{db: db}
	//Set up a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//To close the connection at the end
	defer cancel()
	err := app.db.InitDB(ctx)
	if err != nil {
		return storyTellerApp{}, err
	}
	return app, nil
}

func (app storyTellerApp) CreateStory(ctx context.Context, user, title string) (common.Story, error) {
	if user == "" || title == "" {
		return common.Story{}, fmt.Errorf("user or title empty is not valid")
	}
	story := common.Story{User: user, Title: title, Text: "", ImagesPath: ""}
	story, err := app.db.AddStory(ctx, story)
	if err != nil {
		return common.Story{}, err
	}
	return story, nil
}

func (app storyTellerApp) DeleteStory(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user or title empty is not valid")
	}
	return app.db.DeleteStory(ctx, id)
}

func (app storyTellerApp) FindStory(ctx context.Context, id string) (common.Story, error) {
	if id == "" {
		return common.Story{}, fmt.Errorf("user or title empty is not valid")
	}
	return app.db.FindStory(ctx, id)
}

func (app storyTellerApp) UpdateStory(ctx context.Context, id, user, title, text string) (common.Story, error) {
	if user == "" || title == "" {
		return common.Story{}, fmt.Errorf("user or title empty is not valid")
	}

	return app.db.UpdateStory(ctx, id, common.Story{User: user, Title: title, Text: text})
}
func (app storyTellerApp) FindAllStories(ctx context.Context) ([]common.Story, error) {
	return app.db.FindAllStories(ctx)
}
