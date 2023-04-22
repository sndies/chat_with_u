package main

import (
	"flag"
	"fmt"
	"github.com/sndies/chat_with_u/db"
	httpHandler "github.com/sndies/chat_with_u/middleware/ctx_http_handler"
	myLog "github.com/sndies/chat_with_u/middleware/log"
	"github.com/sndies/chat_with_u/service"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	defer myLog.Flush()

	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	httpHandler.HandleFunc("/", service.IndexHandler)
	httpHandler.HandleFunc("/api/count", service.CounterHandler)
	httpHandler.HandleFunc("/api/wechat_news", service.HandleWechatNews)

	log.Fatal(http.ListenAndServe(":80", nil))
}
