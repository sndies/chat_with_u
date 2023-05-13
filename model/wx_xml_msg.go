package model

import (
	"context"
	"encoding/xml"
	"github.com/sndies/chat_with_u/middleware/log"
	"io"
	"net/http"
	"time"
)

type Msg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	Content      string   `xml:"Content"`
	Recognition  string   `xml:"Recognition"`

	MsgId int64 `xml:"MsgId,omitempty"`
}

func NewMsg(ctx context.Context, r *http.Request) (*Msg, error) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "[NewMsg] io read req err: %v", err)
		return nil, err
	}
	var msg Msg
	if err := xml.Unmarshal(bs, &msg); err != nil {
		log.Errorf(ctx, "[NewMsg] xml unmarshal err: %v", err)
		return nil, err
	}
	return &msg, nil
}

func (msg *Msg) GenerateEchoData(ctx context.Context, s string) []byte {
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}
	bs, err := xml.Marshal(&data)
	if err != nil {
		log.Errorf(ctx, "[GenerateEchoData] xml marshal err: %v", err)
	}
	return bs
}
