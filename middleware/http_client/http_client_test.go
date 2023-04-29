package http_client

import (
	"context"
	"github.com/sndies/chat_with_u/model"
	"testing"
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
		httpProxy = "http://149.28.192.250:3128"
	)

	resBody, err := HttpPost(context.Background(), httpUrl, httpProxy, reqBody, headers)
	if err != nil {
		t.Fatalf("[HttpPost] err: %v", err)
	}
	t.Logf("[HttpPost] resBody: %s", string(resBody))
}
