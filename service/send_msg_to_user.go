package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/sndies/chat_with_u/middleware/http_client"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
)

func SendMsgToUser(ctx context.Context, openId, content string) error {
	// 获取access_token
	accessToken, err := GetAccessToken(ctx)
	if err != nil {
		log.Infof(ctx, "[SendMsgToUser] getAccessToken err: %v", err)
		return err
	}

	// 发起http请求
	msg := model.TextMsgReq{
		OpenId:  openId,
		MsgType: "text",
		Msg: model.TextJson{
			Content: content,
		},
	}
	httpUrl := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + accessToken
	resBytes, err := http_client.HttpPost(ctx, httpUrl, "", msg, nil)
	if err != nil {
		return err
	}

	// 反解析响应
	respBody := model.TextMsgResp{}
	if err := utils.UnMarshal(resBytes, &respBody); err != nil {
		log.Errorf(ctx, "[SendMsgToUser] json decode resp err: %v", err)
		return err
	}

	log.Infof(ctx, "[SendMsgToUser] resp: %s", utils.ToJsonString(respBody))
	if respBody.ErrorCode != 0 {
		return errors.New(strconv.FormatInt(respBody.ErrorCode, 10) + respBody.ErrorMsg)
	}
	return nil
}
