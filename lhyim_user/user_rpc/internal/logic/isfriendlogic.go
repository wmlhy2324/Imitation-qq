package logic

import (
	"context"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_rpc/internal/svc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFriendLogic) IsFriend(in *user_rpc.IsFriendRequest) (*user_rpc.IsFriendResponse, error) {
	res := new(user_rpc.IsFriendResponse)
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, uint(in.User1), uint(in.User2)) {
		res.IsFriend = true
		return res, nil
	}

	return res, nil
}
