package logic

import (
	"context"
	"strconv"

	"lhyim_server/lhyim_user/user_rpc/internal/svc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOnlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOnlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOnlineListLogic {
	return &UserOnlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOnlineListLogic) UserOnlineList(in *user_rpc.UserOnlineRequest) (resp *user_rpc.UserOnlineResponse, err error) {

	resp = new(user_rpc.UserOnlineResponse)
	onlineMap := l.svcCtx.Redis.HGetAll("online").Val()
	for key, _ := range onlineMap {
		val, err1 := strconv.Atoi(key)
		if err1 != nil {
			logx.Error(err1)
			continue
		}
		resp.UserIdList = append(resp.UserIdList, uint32(val))
	}

	return
}
