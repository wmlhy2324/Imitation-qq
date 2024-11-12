package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/common/response"
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		var upGrader = websocket.Upgrader{

			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				logx.Error(err)
				break
			}
			fmt.Println(string(p))
			//发送消息
			conn.WriteMessage(websocket.TextMessage, []byte("xxx"))
		}
		defer conn.Close()
	}
}
