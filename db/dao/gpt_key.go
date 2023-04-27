package dao

import (
	"context"
	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/model"
	"github.com/sndies/chat_with_u/middleware/log"
)

const tableNameGptKey = "gpt_key"

// GetGptKey 查询GptKey
func GetGptKey(ctx context.Context, name string) (*model.GptKey, error) {
	key := new(model.GptKey)

	err := db.Get().Table(tableNameGptKey).Where("name = ?", name).First(key).Error
	if err != nil {
		log.Errorf(ctx, "[GetGptKey] db err: %v", err)
	}

	return key, err
}
