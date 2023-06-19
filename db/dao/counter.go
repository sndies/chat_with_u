package dao

import (
	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/db/db_model"
)

const tableNameCounter = "Counters"

// ClearCounter 清除Counter
func ClearCounter(id int32) error {
	cli := db.Get()
	return cli.Table(tableNameCounter).Delete(&db_model.CounterModel{Id: id}).Error
}

// UpsertCounter 更新/写入counter
func UpsertCounter(counter *db_model.CounterModel) error {
	cli := db.Get()
	return cli.Table(tableNameCounter).Save(counter).Error
}

// GetCounter 查询Counter
func GetCounter(id int32) (*db_model.CounterModel, error) {
	var err error
	var counter = new(db_model.CounterModel)

	cli := db.Get()
	err = cli.Table(tableNameCounter).Where("id = ?", id).First(counter).Error

	return counter, err
}