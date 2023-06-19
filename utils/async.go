package utils

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/log"
	"runtime/debug"
)

func SafeGo(ctx context.Context, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf(ctx, "panic: %s", string(debug.Stack()))
			}
		}()

		f()
	}()
}
