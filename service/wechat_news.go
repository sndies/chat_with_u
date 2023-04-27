package service

import (
	"context"
	"fmt"
	"github.com/sndies/chat_with_u/middleware/log"
	gpt "github.com/sndies/chat_with_u/middleware/gpt_handler"
    "github.com/sndies/chat_with_u/db/model"
	"github.com/sndies/chat_with_u/utils"
	"net/http"
	"time"
)

type NewsReq struct {
	ToUserName   string
	FromUserName string
	MsgType      string
	Content      string
	CreateTime   int64
	MsgId        int64
}

func HandleWechatNews(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	log.Infof(ctx, "receive wechat news, raw_r: %+v", r)

	reqJson := NewsReq{}
	if err := utils.Decoder(r.Body).Decode(&reqJson); err != nil {
		log.Infof(ctx,"json decode req err: %v", err)
		_, _ = fmt.Fprint(w, "内部错误")
		return
	}
	log.Infof(ctx, "receive json req: %s", utils.ToJsonString(reqJson))

	model := model.OpenaiModel{}
	model.Id = "gpt-4"
	answer, err := gpt.Completions(ctx, reqJson.Content, model)
	if err != nil {
		answer = "请稍后再试"
	}
	
	res = map[string]interface{}{
		"ToUserName":   reqJson.FromUserName,
		"FromUserName": reqJson.ToUserName,
		"CreateTime":   time.Now().Unix(),
		"MsgType":      "text",
		"Content":      answer,
	}

	msg, err := utils.Marshal(res)
	if err != nil {
		log.Infof(ctx, "unmarshal res err: %v", err)
		_, _ = fmt.Fprint(w, "内部错误")
		return
	}
	log.Infof(ctx, "response json: %s", utils.ToJsonString(res))

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(msg)
}
