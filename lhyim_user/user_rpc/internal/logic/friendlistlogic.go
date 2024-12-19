package logic

import (
	"context"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_rpc/internal/svc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *user_rpc.FriendListRequest) (*user_rpc.FriendListResponse, error) {
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{

			Limit: -1,
		},
		Where:   l.svcCtx.DB.Where("send_user_id = ? or recv_user_id = ?", in.UserId, in.UserId),
		Preload: []string{"SendUserModel", "RecvUserModel"},
		Debug:   true,
	})
	var list []*user_rpc.FriendInfo
	for _, friend := range friends {
		info := user_rpc.FriendInfo{}
		if friend.SendUserID == uint(in.UserId) {
			//我是发起方
			info = user_rpc.FriendInfo{
				UserId:   uint32(friend.RecvUserID),
				NickName: friend.RecvUserModel.Nickname,
				Avatar:   friend.RecvUserModel.Avatar,
			}
		}
		if friend.RecvUserID == uint(in.UserId) {
			//我是接收方
			info = user_rpc.FriendInfo{

				UserId:   uint32(friend.SendUserID),
				NickName: friend.SendUserModel.Nickname,
				Avatar:   friend.SendUserModel.Avatar,
			}
		}
		list = append(list, &info)
	}
	return &user_rpc.FriendListResponse{FriendList: list}, nil
}
