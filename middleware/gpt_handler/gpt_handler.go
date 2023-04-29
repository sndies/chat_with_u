package gpt_handler

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/http_client"
	"github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/model"
	"github.com/sndies/chat_with_u/utils"
	"os"
)

const gptUseName = "lsndies"
const BASEURL = "https://api.openai.com/v1/"

func Completions(ctx context.Context, msg string, m *model.OpenaiModel) (string, error) {
	if msg == "" {
		log.Infof(ctx, "[Completions] empty msg, return")
		return "", nil
	}
	mod := "gpt-3.5-turbo"
	if m != nil {
		mod = m.Id
	}

	// 调用api
	requestBody := model.OpenaiRequestBody{
		Model: mod,
		Messages: []model.OpenapiRequestMessageItem{
			{Role: "user", Content: msg},
		},
		MaxTokens:        1024,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	proxyUrl := os.Getenv("proxy_url")
	resByte, err := http_client.HttpPost(ctx, BASEURL+"chat/completions", proxyUrl, requestBody, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv("gpt_key"),
	})
	if err != nil {
		return "", err
	}

	// 反序列化为具体结构
	gptResponseBody := model.OpenaiResponseBody{}
	if err := utils.UnMarshal(resByte, gptResponseBody); err != nil {
		log.Errorf(ctx, "[Completions] json decode resp err: %v", err)
		return "", err
	}
	log.Infof(ctx, "[Completions] rawRes: %+v, gptRes: %s", string(resByte), utils.ToJsonString(gptResponseBody))

	// 取回复
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}
	log.Infof(ctx, "[Completions] gpt response text: %s", reply)

	return reply, nil
}
