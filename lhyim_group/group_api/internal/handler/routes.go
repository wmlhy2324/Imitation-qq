// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"lhyim_server/lhyim_group/group_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/group/group",
				Handler: groupCreateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/group/group",
				Handler: groupUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/group/group/:id",
				Handler: groupInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/group/group/:id",
				Handler: groupRemoveHandler(serverCtx),
			},
		},
	)
}
