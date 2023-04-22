package main

import (
	"flag"
	"fmt"
	"github.com/sndies/chat_with_u/db"
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

	http.HandleFunc("/", service.IndexHandler)                    // 测试的
	http.HandleFunc("/api/count", service.CounterHandler)         // 微信云托管自带的,当做以后一个参考
	http.HandleFunc("/api/wechat_news", service.HandleWechatNews) // 自己写的,先固定自动回复一下

	log.Fatal(http.ListenAndServe(":80", nil))
}
