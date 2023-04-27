package ctx_http_handler

import (
	"context"
	"github.com/sndies/chat_with_u/middleware/id_generator"
	"net/http"
)

func HandleFunc(pattern string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(
		pattern,
		func(w http.ResponseWriter, r *http.Request) {
			var (
				ctx      = context.Background()
				logId, _ = id_generator.GenIdInt(ctx)
			)
			ctx = context.WithValue(ctx, "logID", logId)
			handler(ctx, w, r)
		},
	)
}
