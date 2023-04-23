package ctx_http_handler

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
)

func HandleFunc(pattern string, handler func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc(
		pattern,
		func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(context.Background(), "logID", uuid.NewV1())
			log.Printf("-------> logID: %v", ctx.Value("logID"))
			handler(ctx, w, r)
		},
	)
}
