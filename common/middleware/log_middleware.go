package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lhyim_server/common/log_stash"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
	Body []byte
}

func (w *Writer) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}
func LogActionMiddleware(pusher *log_stash.Pusher) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//设置入参
			pusher.SetRequest(r)
			pusher.SetHeader(r)

			clientIP := httpx.GetRemoteAddr(r)
			ctx := context.WithValue(r.Context(), "userID", r.Header.Get("User-ID"))
			ctx = context.WithValue(ctx, "clientIP", clientIP)
			var nw = Writer{
				ResponseWriter: w,
			}
			next(&nw, r.WithContext(ctx))
			//设置响应
			if pusher.GetResponse() {
				//读响应体
				pusher.SetResponse(string(nw.Body))
			}

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
