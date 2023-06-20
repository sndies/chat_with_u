package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sndies/chat_with_u/middleware/gpt_handler"
	"github.com/sndies/chat_with_u/middleware/log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sndies/chat_with_u/db/dao"
	"github.com/sndies/chat_with_u/db/db_model"
	"gorm.io/gorm"
)

// JsonResult 返回结构
type JsonResult struct {
	Status   int         `json:"status"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// IndexHandler 计数器接口
func IndexHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Infof(ctx, "index request: %+v", *r)
	echoStr := r.URL.Query().Get("echostr")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(echoStr))
	//data, err := getIndex()
	//if err != nil {
	//	fmt.Fprint(w, "内部错误")
	//	return
	//}
	//fmt.Fprint(w, data)
}

// CounterHandler 计数器接口
func CounterHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}

	log.Infof(ctx, "http request: %+v", r)
	if r.Method == http.MethodGet {
		reply, err := gpt_handler.Completions(ctx, "今天星期几", nil)
		if err != nil {
			reply = "出现错误"
			res.ErrorMsg = err.Error()
		}
		res.Data = reply
		w.WriteHeader(200)
		//counter, err := getCurrentCounter()
		//if err != nil {
		//	res.Status = -1
		//	res.ErrorMsg = err.Error()
		//} else {
		//	res.Data = counter.Count
		//}
	} else if r.Method == http.MethodPost {
		count, err := modifyCounter(r)
		if err != nil {
			res.Status = -1
			res.ErrorMsg = err.Error()
		} else {
			res.Data = count
		}
	} else {
		res.Status = -1
		res.ErrorMsg = fmt.Sprintf("请求方法 %s 不支持", r.Method)
	}

	msg, err := json.Marshal(res)
	if err != nil {
		fmt.Fprint(w, "内部错误")
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(msg)
}

// modifyCounter 更新计数，自增或者清零
func modifyCounter(r *http.Request) (int32, error) {
	action, err := getAction(r)
	if err != nil {
		return 0, err
	}

	var count int32
	if action == "inc" {
		count, err = upsertCounter(r)
		if err != nil {
			return 0, err
		}
	} else if action == "clear" {
		err = clearCounter()
		if err != nil {
			return 0, err
		}
		count = 0
	} else {
		err = fmt.Errorf("参数 action : %s 错误", action)
	}

	return count, err
}

// upsertCounter 更新或修改计数器
func upsertCounter(r *http.Request) (int32, error) {
	currentCounter, err := getCurrentCounter()
	var count int32
	createdAt := time.Now()
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		count = 1
		createdAt = time.Now()
	} else {
		count = currentCounter.Count + 1
		createdAt = currentCounter.CreatedAt
	}

	counter := &db_model.CounterModel{
		Id:        1,
		Count:     count,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
	err = dao.UpsertCounter(counter)
	if err != nil {
		return 0, err
	}
	return counter.Count, nil
}

func clearCounter() error {
	return dao.ClearCounter(1)
}

// getCurrentCounter 查询当前计数器
func getCurrentCounter() (*db_model.CounterModel, error) {
	counter, err := dao.GetCounter(1)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// getAction 获取action
func getAction(r *http.Request) (string, error) {
	decoder := json.NewDecoder(r.Body)
	body := make(map[string]interface{})
	if err := decoder.Decode(&body); err != nil {
		return "", err
	}
	defer r.Body.Close()

	action, ok := body["action"]
	if !ok {
		return "", fmt.Errorf("缺少 action 参数")
	}

	return action.(string), nil
}

// getIndex 获取主页
func getIndex() (string, error) {
	b, err := ioutil.ReadFile("./index.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
