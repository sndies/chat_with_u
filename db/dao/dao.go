package dao

import (
	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/model"
)

const tableName = "Counters"

// ClearCounter 清除Counter
func ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableName).Delete(&model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func UpsertCounter(counter *model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableName).Save(counter).Error
}

// GetCounter 查询Counter
func GetCounter(id int32) (*model.CounterModel, error) {
	var err error
	var counter = new(model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableName).Where("id = ?", id).First(counter).Error

	return counter, err
}