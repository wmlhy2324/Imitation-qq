package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lhyim_server/common/log_stash"
	"net/http"
)

func LogActionMiddleware(pusher *log_stash.Pusher) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//设置入参
			pusher.SetRequest(r)

			clientIP := httpx.GetRemoteAddr(r)
			ctx := context.WithValue(r.Context(), "userID", r.Header.Get("User-ID"))
			ctx = context.WithValue(ctx, "clientIP", clientIP)
			next(w, r.WithContext(ctx))
			//设置响应
		}
	}
}
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		clientIP := httpx.GetRemoteAddr(r)
		//ctx := context.WithValue(r.Context(), "clientIP", clientIP)
		//ctx = context.WithValue(ctx, "userID", r.Header.Get("User-ID"))
		ctx := context.WithValue(r.Context(), "userID", r.Header.Get("User-ID"))
		ctx = context.WithValue(ctx, "clientIP", clientIP)
		next(w, r.WithContext(ctx))
	}
}
