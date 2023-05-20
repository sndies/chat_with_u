package service

import (
	"context"
	"github.com/sndies/chat_with_u/consts"
	"github.com/sndies/chat_with_u/middleware/cache"
	"github.com/sndies/chat_with_u/middleware/gpt_handler"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
	"net/http"
	"strconv"
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
	log.Infof(ctx, "receive wechat news, raw_r: %+v", r)

	// 这段是用来接入微信开发者验证的
	if r.Method == http.MethodGet {
		echoStr := r.URL.Query().Get("echostr")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(echoStr))
		return
	}

	// 从xml解析为结构体
	reqJson, err := model.NewMsg(ctx, r)
	if err != nil {
		echo(w, []byte("success"))
		return
	}
	log.Infof(ctx, "receive json req: %s", utils.ToJsonString(reqJson))

	// 调用处理逻辑
	reply := queryAndWrapRes(ctx, reqJson)

	// 返回结果
	if reply != "" {
		echo(w, reqJson.GenerateEchoData(ctx, reply))
	}
}

func queryAndWrapRes(ctx context.Context, req *model.Msg) (reply string) {
	// 出口日志
	start := time.Now()
	defer log.Infof(ctx, "[queryAndWrapRes] req: %s, reply: %s, cost: %v", utils.ToJsonString(req), reply, time.Since(start))

	// 检查入参
	req.Content = strings.TrimSpace(req.Content)
	if pass, errMsg := checkReqMsg(ctx, req.Content); !pass {
		return errMsg
	}

	// 同一用户对话的串行控制,先用内存吧,后续迁移到mysql
	_, ok := cache.Get(ctx, req.FromUserName)
	if ok {
		reply = "上个问题正在处理中，请稍等..."
		return
	}

	// 将user放入缓存
	_ = cache.Add(ctx, req.FromUserName, true, time.Second*5)
	defer func() {
		cache.Del(ctx, req.FromUserName)
	}()

	// 调用gpt
	reply = invokeCompletion(ctx, req)

	return
}

func invokeCompletion(ctx context.Context, req *model.Msg) string {
	var (
		ch    chan string
		msgId = strconv.FormatInt(req.MsgId, 10)
	)

	v, ok := cache.Get(ctx, msgId)
	if !ok {
		ch = make(chan string)
		_ = cache.Add(ctx, msgId, ch, time.Minute)
		// 发起请求
		reply, err := gpt_handler.Completions(ctx, req.Content, nil)
		if err != nil {
			ch <- err.Error()
		}
		// 出错
		if len(reply) == 0 {
			ch <- "openai请求超时"
		}
		ch <- reply
	} else {
		ch = v.(chan string)
	}

	select {
	case result := <-ch:
		cache.Del(ctx, msgId)
		return result
	case <-time.After(time.Second * 5):
		// 超时不要回答，会重试的
	}

	return ""
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

func echo(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
