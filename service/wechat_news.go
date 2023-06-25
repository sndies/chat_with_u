package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/sndies/chat_with_u/consts"
	"github.com/sndies/chat_with_u/db/dao"
	"github.com/sndies/chat_with_u/db/db_model"
	"github.com/sndies/chat_with_u/middleware/cache"
	"github.com/sndies/chat_with_u/middleware/gpt_handler"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
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

	// === 特殊: 这段是用来接入微信开发者验证的 ===
	if r.Method == http.MethodGet {
		echoStr := r.URL.Query().Get("echostr")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(echoStr))
		return
	}

	// 1.从xml解析为结构体
	reqJson, err := model.NewMsg(ctx, r)
	if err != nil {
		echo(w, []byte("success"))
		return
	}
	log.Infof(ctx, "receive json req: %s", utils.ToJsonString(reqJson))

	// 2.调用处理逻辑
	reply := queryAndWrapRes(ctx, reqJson)

	// 3.返回结果
	if reply != "" {
		echo(w, reqJson.GenerateEchoData(ctx, reply))
	}
}

func queryAndWrapRes(ctx context.Context, req *model.Msg) (reply string) {
	// 出口日志
	start := time.Now()
	defer log.Infof(ctx, "[queryAndWrapRes] req: %s, reply: %s, cost: %v", utils.ToJsonString(req), reply, time.Since(start))

	// 1.检查入参
	req.Content = strings.TrimSpace(req.Content)
	if pass, errMsg := checkReqMsg(ctx, req.Content); !pass {
		return errMsg
	}

	// 2.同一用户对话的串行控制,暂时先用内存,后续可以迁移到mysql
	_, ok := cache.Get(ctx, req.FromUserName)
	if ok {
		reply = "上个问题正在处理中，请稍等..."
		return
	}

	// 3.将user放入缓存
	_ = cache.Add(ctx, req.FromUserName, true, time.Second*10)

	// 4.用户当日免费额度检查
	isRunOut, _ := IsUserHasRunOutOfQuota(ctx, req.FromUserName)
	if isRunOut {
		reply = "抱歉，每天仅有2次免费问答机会"
		return
	}

	// 5.用msgId查询数据库
	qna, err := dao.GetGptQNAByMsgId(ctx, req.MsgId)
	if qna != nil {
		// 4.1 如果已经有答案直接返回
		if qna.Answer != "" {
			reply = qna.Answer
			return
		}
		// 4.2 没有答案,告诉用户正在处理中
		reply = "收到，请稍等..."
		return

	} else {
		// 4.3 如果没有答案, 写记录，先返回处理中
		err = dao.InsertGptQNA(ctx, &db_model.GptQNA{
			MsgId:        req.MsgId,
			FromUserName: req.FromUserName,
			Question:     req.Content,
			CreatedAt:    time.Now(),
		})
		if err != nil {
			reply = "网络不稳定，请再发送一次"
			return
		}
		reply = "收到，请稍等..."
	}

	// 5.异步请求gpt
	utils.SafeGo(ctx, func() {
		log.Infof(ctx, "[queryAndWrapRes] 进入异步流程")
		defer cache.Del(ctx, req.FromUserName) // 删除缓存,用户可以发起下一次请求
		// 5.1 发起请求
		reply, err = gpt_handler.Completions(ctx, req.Content, nil)
		if err != nil || len(reply) == 0 {
			return
		}
		// 5.2 写进数据库
		err = dao.UpdateAnswerByMsgId(ctx, reply, req.MsgId)
		if err != nil {
			return
		}
		log.Infof(ctx, "[queryAndWrapRes] 写数据库成功")
		// 5.3 将结果主动推送给用户
		_ = SendMsgToUser(ctx, req.FromUserName, reply)

	})

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

func echo(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
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
		//log.Infof(ctx, "[invokeCompletion] channel timeout, req: %s", utils.ToJsonString(req))
		glog.Infof("[invokeCompletion] channel timeout, req: %s", utils.ToJsonString(req))
	}

	return ""
}
