package service

import (
	"context"
	"fmt"
	"github.com/sndies/chat_with_u/consts"
	"github.com/sndies/chat_with_u/middleware/cache"
	"github.com/sndies/chat_with_u/middleware/gpt_handler"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/utils"
	"net/http"
	"strings"
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
		log.Infof(ctx, "json decode req err: %v", err)
		_, _ = fmt.Fprint(w, "内部错误")
		return
	}
	log.Infof(ctx, "receive json req: %s", utils.ToJsonString(reqJson))

	reply := queryAndWrapRes(ctx, reqJson.FromUserName, reqJson.Content, 1*time.Minute)
	res = map[string]interface{}{
		"ToUserName":   reqJson.FromUserName,
		"FromUserName": reqJson.ToUserName,
		"CreateTime":   time.Now().Unix(),
		"MsgType":      "text",
		"Content":      reply,
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

func queryAndWrapRes(ctx context.Context, uid, msg string, timeout time.Duration) (reply string) {
	// 出口日志
	start := time.Now()
	defer log.Infof(ctx, "[queryAndWrapRes] uid: %s, msg: %s, reply: %s, cost: %v", uid, msg, &reply, time.Since(start).Milliseconds())

	// 超时设置
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 检查入参
	msg = strings.TrimSpace(msg)
	if pass, errMsg := checkReqMsg(ctx, msg); !pass {
		return errMsg
	}

	// 同一用户对话的串行控制,先用内存吧,后续迁移到mysql
	_, ok := cache.Get(ctx, uid)
	if ok {
		reply = "上个问题正在处理中，请稍等..."
		return
	}

	_ = cache.Add(ctx, uid, true, time.Minute*2)
	defer func() {
		cache.Del(ctx, uid)
	}()

	// 发起请求
	reply, err := gpt_handler.Completions(ctx, msg, nil)
	if err != nil {
		return err.Error()
	}

	// 超时结束
	var done bool
	for !done {
		select {
		case <-ctx.Done():
			done = true
		default:
			done = reply != ""
		}
	}

	// 出错
	if len(reply) == 0 {
		reply = "openai请求超时"
		return
	}

	return
}

func checkReqMsg(ctx context.Context, msg string) (bool, string) {
	length := len([]rune(msg))

	if length <= 1 {
		return false, "请说详细些..."
	}
	if length > consts.MaxQuestionLength {
		return false, "问题字数超出设定限制，请精简问题"
	}

	return true, ""
}
