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
type TextMsgReq struct {
	OpenId  string   `json:"touser"`
	MsgType string   `json:"msgtype"`
	Msg     TextJson `json:"text"`
}
type TextMsgResp struct {
	ErrorCode int64  `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

type GetAccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	AppId        string `json:"appid"`
	Secret       string `json:"secret"`
	ForceRefresh bool `json:"force_refresh"`
}
type GetAccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int64  `json:"expires_in"`
}
