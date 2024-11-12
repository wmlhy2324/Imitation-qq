package handler

import (
	"lhyim_server/common/response"
	"lhyim_server/lhyim_chat/chat_api/internal/logic"
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func chatDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewChatDeleteLogic(r.Context(), svcCtx)
		resp, err := l.ChatDelete(&req)
		response.Response(r, w, resp, err)

	}
}
