package memorydb

import (
	"context"

	"github.com/idalmasso/storytellers/backend/common"
)

type MemoryDatabase struct {
	stories map[string]common.Story
}

func (db *MemoryDatabase) InitDB(ctx context.Context) error {
	db.stories = make(map[string]common.Story, 0)
	return nil
}
