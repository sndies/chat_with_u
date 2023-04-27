package http_client

import (
	"bytes"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHttpPost(t *testing.T) {
	var (
		reqBody = model.OpenaiRequestBody{
			Model: "gpt-3.5-turbo",
			Messages: []model.OpenapiRequestMessageItem{
				{Role: "user", Content: "今天星期几"},
			},
			MaxTokens:        1024,
			Temperature:      0.7,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		}
		headers = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + "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",
		}
		httpUrl   = "https://api.openai.com/v1/chat/completions"
		httpProxy = ""
	)

	requestData, err := utils.Marshal(reqBody)
	if err != nil {
		t.Fatalf("[HttpPost] unmarshal reqBody err: %v", err)
	}
	client := &http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("POST", httpUrl, bytes.NewReader(requestData))
	if err != nil {
		t.Fatalf("[HttpPost] new http request err: %v", err)
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
		proxyURL, err := url.Parse(httpProxy)
		if err != nil {
			t.Fatalf("[HttpPost] parse err: %v", err)
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// 调用
	response, err := client.Do(req)
	t.Logf("[HttpPost] req: %+v, \nres: %+v", req, response)
	if err != nil {
		t.Fatalf("[HttpPost] http post err: %v", err)
	}

	// 读取响应
	defer func() { _ = response.Body.Close() }()
	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("[HttpPost] read body err: %v", err)
	}

	t.Logf("[HttpPost] resBody: %s", string(resBody))
}
