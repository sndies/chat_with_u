package model

import "time"

// GptKey 存储api_key
type GptKey struct {
	Id        int32     `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Key       string    `gorm:"column:key" json:"key"`
	CreatedAt time.Time `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"updatedAt"`
}
