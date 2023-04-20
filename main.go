package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/sndies/chat_with_u/db"
	"github.com/sndies/chat_with_u/service"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
