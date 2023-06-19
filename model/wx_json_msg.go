package model

/*
{
    "touser":"OPENID",
    "msgtype":"text",
    "text":
    {
         "content":"Hello World"
    }
}
*/
type TextJson struct {
	Content string `json:"content"`
}
type TextMsgJson struct {
	OpenId  string   `json:"touser"`
	MsgType string   `json:"msgtype"`
	Msg     TextJson `json:"text"`
}
