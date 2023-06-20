package db_model

import (
	"time"
)

type GptQNA struct {
	Id           int32     `gorm:"column:id" json:"id"`
	MsgId        int64     `gorm:"column:msg_id" json:"msg_id"`
	FromUserName string    `gorm:"column:from_user_name" json:"from_user_name"`
	Question     string    `gorm:"column:question" json:"question"`
	Answer       string    `gorm:"column:answer" json:"answer"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}
