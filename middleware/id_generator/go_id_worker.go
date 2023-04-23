package id_generator

import (
	"context"
	"fmt"
	idworker "github.com/gitstliu/go-id-worker"
	"github.com/sndies/chat_with_u/middleware/log"
)

var currWorker *idworker.IdWorker

func Init() {
	currWorker = &idworker.IdWorker{}
	if err := currWorker.InitIdWorker(1000, 1); err != nil {
		panic(fmt.Sprintf("idWorker init failed with %+v", err))
	}
}

func GenIdInt(ctx context.Context) (int64, error) {
	id, err := currWorker.NextId()
	if err != nil {
		log.Errorf(ctx, "currWorker.nextId() err: %v", err)
	}
	return id, err
}
