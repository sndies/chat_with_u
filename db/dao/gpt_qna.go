package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"

	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/db_model"
	"github.com/sndies/chat_with_u/middleware/log"
)

const tableNameGptQna = "gpt_qna"

func GetGptQNAByMsgId(ctx context.Context, msgId int64) (*db_model.GptQNA, error) {
	qna := new(db_model.GptQNA)

	err := db.Get().Table(tableNameGptQna).Where("msgId = ?", msgId).First(qna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Errorf(ctx, "[GetGptQNAByMsgId] db err: %v", err)
	}

	return qna, err
}

func InsertGptQNA(ctx context.Context, qna *db_model.GptQNA) error {
	err := db.Get().Table(tableNameGptQna).Create(qna).Error
	if err != nil {
		log.Errorf(ctx, "[InsertGptQNA] db err: %v", err)
	}
	return err
}

func UpdateAnswerByMsgId(ctx context.Context, answer string, msgId int64) error {
	err := db.Get().Table(tableNameGptQna).
		Where("msg_id = ?", msgId).
		UpdateColumns(map[string]interface{}{
			"answer": answer,
		}).
		Error
	if err != nil {
		log.Errorf(ctx, "[UpdateAnswerByMsgId] db err: %v", err)
	}
	return err
}
