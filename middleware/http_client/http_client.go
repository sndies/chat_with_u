package http_client

import (
	"bytes"
	"context"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func HttpPost(ctx context.Context, httpUrl, httpProxy string, reqBody interface{}, headers map[string]string) ([]byte, error) {
	// req处理
	requestData, err := utils.Marshal(reqBody)
	if err != nil {
		log.Errorf(ctx, "[HttpPost] marshal reqBody err: %v", err)
		return nil, err
	}
	client := &http.Client{Timeout: time.Second * 200}
	req, err := http.NewRequest("POST", httpUrl, bytes.NewReader(requestData))
	if err != nil {
		log.Errorf(ctx, "[HttpPost] new http request err: %v", err)
		return nil, err
	}

	// header处理
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 代理
	if httpProxy != "" {
		proxyURL, _ := url.Parse(httpProxy)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// 调用
	response, err := client.Do(req)
	log.Infof(ctx, "[HttpPost] req: %+v, res: %+v", req, response)
	if err != nil {
		log.Errorf(ctx, "[HttpPost] http post err: %v", err)
		return nil, err
	}

	// 读取响应
	defer func() { _ = response.Body.Close() }()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf(ctx, "[HttpPost] read body err: %v", err)
		return nil, err
	}

	return body, nil
}
