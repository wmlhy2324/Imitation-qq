package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := httpx.GetRemoteAddr(r)
		ctx := context.WithValue(r.Context(), "clientIP", clientIP)
		ctx = context.WithValue(ctx, "User-ID", r.Header.Get("User-ID"))
		next(w, r.WithContext(ctx))
	}
}
