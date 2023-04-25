package gtp_handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"context"
	"github.com/sndies/chat_with_u/db/model"
	"github.com/sndies/chat_with_u/utils"
	"github.com/sndies/chat_with_u/middleware/log"
)

const BASEURL = "https://api.openai.com/v1/"
const APIKEY = ""

func Completions(ctx context.Context, msg string, m model.OpenaiModel) (string, error) {
	if msg == "" {
		return "", nil
	}

	mod := "gpt-4"
	// TODO 怎么判断对象是否为nil
	// if &m != nil {
	// 	mod = m.Id
	// }
	requestBody := model.OpenaiRequestBody{
		Model:            mod,
		Prompt:           msg,
		MaxTokens:        2048,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	requestData, err := utils.Marshal(requestBody)
	if err != nil {
		log.Infof(ctx, "unmarshal requestData err: %v", err)
		_, _ = fmt.Printf("内部错误")
		return "", err
	}
	log.Infof(ctx, "resquest json: %s", requestData)

	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + APIKEY)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := model.OpenaiResponseBody{}
	if err := utils.UnMarshal(body, gptResponseBody); err != nil {
		log.Infof(ctx,"json decode resp err: %v", err)
		_, _ = fmt.Printf("内部错误")
		return "", err
	}
	
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}
	log.Infof(ctx, "gpt response text: %s \n", reply)
	return reply, nil
}