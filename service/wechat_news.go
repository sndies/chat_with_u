package service

import (
	"fmt"
	"github.com/golang/glog"
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

func HandleWechatNews(w http.ResponseWriter, r *http.Request) {
	res := &JsonResult{}
	glog.Infof("receive wechat news, raw_r: %+v", r)

	reqJson := NewsReq{}
	if err := utils.Decoder(r.Body).Decode(&reqJson); err != nil {
		glog.Infof("json decode req err: %v", err)
		_, _ = fmt.Fprint(w, "内部错误")
		return
	}
	glog.Infof("receive json req: %s", utils.ToJsonString(reqJson))

	res.Status = 200
	res.Data = map[string]interface{}{
		"ToUserName":   reqJson.FromUserName,
		"FromUserName": reqJson.ToUserName,
		"CreateTime":   time.Now().Unix(),
		"MsgType":      "text",
		"Content":      "收到,等我回复吧",
	}

	msg, err := utils.Marshal(res)
	if err != nil {
		glog.Infof("unmarshal res err: %v", err)
		_, _ = fmt.Fprint(w, "内部错误")
		return
	}
	glog.Infof("response: %s", utils.ToJsonString(res))

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(msg)
}
