package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/db_model"
	"github.com/sndies/chat_with_u/middleware/log"
)

const tableNameGptQna = "gpt_qna"

func GetGptQNAByMsgId(ctx context.Context, msgId int64) (*db_model.GptQNA, error) {
	qna := new(db_model.GptQNA)

	err := db.Get().Table(tableNameGptQna).Where("msg_id = ?", msgId).First(qna).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Errorf(ctx, "[GetGptQNAByMsgId] db err: %v", err)
	}

	return qna, err
}

func GetGptQNAByOpenId(ctx context.Context, openId string, startTime, endTime time.Time) ([]*db_model.GptQNA, error) {
	var qnaList []*db_model.GptQNA

	err := db.Get().Table(tableNameGptQna).
		Where("from_user_name = ? and created_at >= ? and created_at < ?", openId, startTime, endTime).
		Find(&qnaList).
		Error
	if err != nil {
		log.Errorf(ctx, "[GetGptQNAByOpenId] db err: %+v", err)
	}

	return qnaList, err
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
			"answer":     answer,
			"updated_at": time.Now(),
		}).
		Error
	if err != nil {
		log.Errorf(ctx, "[UpdateAnswerByMsgId] db err: %v", err)
	}
	return err
}
