package dao

import (
	"context"
	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/db_model"
	"github.com/sndies/chat_with_u/middleware/log"
)

const tableNameGptKey = "gpt_key"

// GetGptKey 查询GptKey
func GetGptKey(ctx context.Context, name string) (*db_model.GptKey, error) {
	key := new(db_model.GptKey)

	err := db.Get().Table(tableNameGptKey).Where("name = ?", name).First(key).Error
	if err != nil {
		log.Errorf(ctx, "[GetGptKey] db err: %v", err)
	}

	return key, err
}
