package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_chat/chat_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTopLogic {
	return &UserTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserTopLogic) UserTop(req *types.UserTopRequest) (resp *types.UserTopResponse, err error) {

	//是否是好友
	res, err := l.svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
		User1: uint32(req.UserID),
		User2: uint32(req.FriendID),
	})
	if err != nil {
		return nil, err
	}
	if !res.IsFriend {
		return nil, errors.New("你们还不是好友")
	}
	var topUser chat_models.TopUserModel
	err = l.svcCtx.DB.Take(&topUser, "user_id = ? and top_user_id = ?", req.UserID, req.FriendID).Error
	if err != nil {
		//没有置顶
		l.svcCtx.DB.Create(&chat_models.TopUserModel{
			UserID:    req.UserID,
			TopUserID: req.FriendID,
		})
		return nil, nil
	}
	//已经有置顶
	l.svcCtx.DB.Model(&chat_models.TopUserModel{}).Delete("user_id = ? and top_user_id = ?", req.UserID, req.FriendID)
	return
}