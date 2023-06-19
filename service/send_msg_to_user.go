package service

//import (
//	"context"
//	"github.com/sndies/chat_with_u/middleware/http_client"
//	"github.com/sndies/chat_with_u/model"
//)
//
//func SendMsgToUser(ctx context.Context, openId, content string) error {
//	msg := model.TextMsgJson{
//		OpenId:  openId,
//		MsgType: "text",
//		Msg: model.TextJson{
//			Content: content,
//		},
//	}
//	httpUrl := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ""
//
//	// 发起http请求
//	resBytes, err := http_client.HttpPost(ctx, "", "", msg, nil)
//	if err != nil {
//		return err
//	}
//}
