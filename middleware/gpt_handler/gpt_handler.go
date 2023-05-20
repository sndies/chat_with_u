package gpt_handler

import (
	"context"
	"github.com/golang/glog"
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
	resByte, err := http_client.HttpPost(ctx, BASEURL+"chat/completions", "", requestBody, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv("gpt_key"),
	})
	if err != nil {
		return "", err
	}
	//log.Infof(ctx, "[Completions] http_res: %s", string(resByte))
	glog.Infof("[Completions] http_res: %s", string(resByte))

	// 反序列化为具体结构
	gptResponseBody := model.OpenaiResponseBody{}
	if err := utils.UnMarshal(resByte, &gptResponseBody); err != nil {
		log.Errorf(ctx, "[Completions] json decode resp err: %v", err)
		return "", err
	}
	//log.Infof(ctx, "[Completions] gptRes: %s", utils.ToJsonString(gptResponseBody))
	glog.Infof("[Completions] gptRes: %s", utils.ToJsonString(gptResponseBody))

	// 取回复
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, choice := range gptResponseBody.Choices {
			reply = choice.Message.Content
			break
		}
	}
	//log.Infof(ctx, "[Completions] gpt response text: %s", reply)
	glog.Infof("[Completions] gpt response text: %s", reply)

	return reply, nil
}
