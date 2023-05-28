package web

import (
	"context"
	"net/http"
)

func Handle(appCtx context.Context, addr string, rootHandler http.Handler) {
	muxer := http.NewServeMux()

	muxer.Handle("/", rootHandler)

	Serve(appCtx, addr, muxer)
}
