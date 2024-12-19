package handler

import (
	"lhyim_server/common/response"
	"lhyim_server/lhyim_group/group_api/internal/logic"
	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func groupMemberAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupMemberAddRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGroupMemberAddLogic(r.Context(), svcCtx)
		resp, err := l.GroupMemberAdd(&req)
		response.Response(r, w, resp, err)

	}
}
