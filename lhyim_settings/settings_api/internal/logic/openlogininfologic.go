package logic

import (
	"context"

	"lhyim_server/lhyim_settings/settings_api/internal/svc"
	"lhyim_server/lhyim_settings/settings_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Open_login_infoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpen_login_infoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Open_login_infoLogic {
	return &Open_login_infoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Open_login_infoLogic) Open_login_info() (resp []types.OpenLoginResponse, err error) {
	for _, s := range l.svcCtx.Config.OpenLoginList {
		resp = append(resp, types.OpenLoginResponse{
			Name: s.Name,
			Ican: s.Icon,
			Href: s.Href,
		})

	}
	return
}
