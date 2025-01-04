package logic

import (
	"context"
	"encoding/json"
	"errors"
	"lhyim_server/lhyim_user/user_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendInfoLogic {
	return &FriendInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendInfoLogic) FriendInfo(req *types.FriendInfoRequest) (resp *types.FriendInfoResponse, err error) {
	var friend user_models.FriendModel
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("非好友")
	}
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoRequest{
		UserId: uint32(req.FriendID),
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	var user user_models.UserModel
	json.Unmarshal(res.Data, &user)
	return &types.FriendInfoResponse{
		UserID:   int64(user.ID),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Abstract: user.Abstract,
		Notice:   friend.GetUserNotice(req.UserID),
	}, nil
}
