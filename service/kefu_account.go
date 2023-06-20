package service

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/http_client"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
)

func AddKefuAccount(ctx context.Context, accessToken ...string) {
	var (
		token string
		err   error
	)
	if len(accessToken) == 0 {
		token, err = GetAccessToken(ctx)
		if err != nil {
			log.Errorf(ctx, "[AddKefuAccount] get accessToken err: %v", err)
			return
		}
	} else {
		token = accessToken[0]
	}

	httpUrl := "https://api.weixin.qq.com/customservice/kfaccount/add?access_token=" + token
	req := map[string]string{
		"kf_account": "ga@GrowthAI",
		"nickname":   "智模型",
		"password":   "helloworld123",
	}
	resBytes, err := http_client.HttpPost(ctx, httpUrl, "", req, nil)
	if err != nil {
		return
	}
	log.Infof(ctx, "[AddKefuAccount] http resp: %s", string(resBytes))

	respBody := model.TextMsgResp{}
	if err := utils.UnMarshal(resBytes, &respBody); err != nil {
		log.Errorf(ctx, "[AddKefuAccount] json decode resp err: %v", err)
		return
	}

	log.Infof(ctx, "[AddKefuAccount] resp: %s", utils.ToJsonString(respBody))
}
