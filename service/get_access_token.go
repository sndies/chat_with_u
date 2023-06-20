package service

import (
	"context"
	"time"

	"github.com/sndies/chat_with_u/consts"
	"github.com/sndies/chat_with_u/middleware/cache"
	"github.com/sndies/chat_with_u/middleware/http_client"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
)

func GetAccessToken(ctx context.Context) (string, error) {
	// 先从缓存中取
	cacheVal, exist := cache.Get(ctx, consts.AccessTokenKey)
	if exist {
		if val, ok := cacheVal.(string); ok {
			log.Infof(ctx, "[GetAccessToken] cache token: %s", val)
			return val, nil
		}
	}

	// 没有缓存开始请求
	httpUrl := "https://api.weixin.qq.com/cgi-bin/stable_token"
	req := model.GetAccessTokenReq{
		GrantType: "client_credential",
		AppId:     "wxf2492d93ee4fb797",
		//Secret:    os.Getenv("app_secret"),
		Secret: "84059f0ee28e8f16c235422dbe6f6d1d",
	}

	// 发起http请求
	resBytes, err := http_client.HttpPost(ctx, httpUrl, "", req, nil)
	if err != nil {
		return "", err
	}
	log.Infof(ctx, "[GetAccessToken] http resp: %s", string(resBytes))

	// 反序列化为具体结构
	responseBody := model.GetAccessTokenResp{}
	if err := utils.UnMarshal(resBytes, &responseBody); err != nil {
		log.Errorf(ctx, "[GetAccessToken] json decode resp err: %v", err)
		return "", err
	}
	log.Infof(ctx, "[GetAccessToken] token: %s", responseBody.AccessToken)

	// 写缓存
	validSeconds := time.Duration(responseBody.ExpireIn) * time.Second
	_ = cache.Add(ctx, consts.AccessTokenKey, responseBody.AccessToken, validSeconds)

	return responseBody.AccessToken, nil
}
