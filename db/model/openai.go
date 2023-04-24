package model

// OpenaiModel openai model
// https://platform.openai.com/docs/api-reference/chat/create
type OpenaiModel struct {
	id        int32 
	object    string
	ownedBy   string `json:"owned_by"`
	permission [] string
}

type MessageModel struct {
	role	string
	content	string
	name	string // optional
}

type OpenaiRequestBody struct {
	model	string // required
	messages [] MessageModel
	suffix	string // optional defaults to null
	max_token int32 // optional defaults to 16
	temperature	int // optional defaults to 1
	top_p	int // optional 
	n	int32 // How many completions to generate for each prompt.
	stream	bool // Whether to stream back partial progress default false
}