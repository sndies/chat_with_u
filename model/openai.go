package model

// OpenaiModel openai model
// https://platform.openai.com/docs/api-reference/chat/create
type OpenaiModel struct {
	Id         string
	Object     string
	OwnedBy    string `json:"owned_by"`
	Permission [] string
}

type MessageModel struct {
	Role    string
	Content string
	Name    string // optional
}

type OpenapiRequestMessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenaiRequestBody struct {
	Model            string                      `json:"model"`
	Messages         []OpenapiRequestMessageItem `json:"messages"`
	MaxTokens        int                         `json:"max_tokens"`
	Temperature      float32                     `json:"temperature"`
	TopP             int                         `json:"top_p"`
	FrequencyPenalty int                         `json:"frequency_penalty"`
	PresencePenalty  int                         `json:"presence_penalty"`
}

type OpenaiResponseBody struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int                      `json:"created"`
	Model   string                   `json:"model"`
	Choices []map[string]interface{} `json:"choices"`
	Usage   map[string]interface{}   `json:"usage"`
}
