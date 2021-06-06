package types

import (
	"context"

	"github.com/idalmasso/storytellers/backend/common"
)

type StoryTellerAppInterface interface{
	CreateStory(ctx context.Context, user, title string) (common.Story, error) 
	DeleteStory(ctx context.Context, id string) error
	FindStory(ctx context.Context, id string) (common.Story, error)
	UpdateStory(ctx context.Context, id, user, title, text string) (common.Story, error)
	FindAllStories(ctx context.Context) ([]common.Story, error)
}
