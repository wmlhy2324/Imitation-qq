package logic

import (
	"context"
	"fmt"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListRequest) (resp *types.FriendListResponse, err error) {
	//var count int64
	//l.svcCtx.DB.Model(&user_models.FriendModel{}).Where("send_user_id = ? or recv_user_id = ?", req.UserID, req.UserID).Count(&count)
	//var friends []user_models.FriendModel
	//
	//if req.Limit <= 0 {
	//	req.Limit = 10
	//}
	//if req.Page <= 0 {
	//	req.Page = 1
	//
	//}
	//offset := (req.Page - 1) * req.Limit
	//l.svcCtx.DB.Preload("SendUserModel").Preload("RecvUserModel").Limit(req.Limit).Offset(offset).Find(&friends, "send_user_id = ? or recv_user_id = ?", req.UserID, req.UserID)

	var list []types.FriendInfoResponse
	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"SendUserModel", "RecvUserModel"},
	})
	fmt.Println(friends)

	for _, friend := range friends {
		info := types.FriendInfoResponse{}
		if friend.SendUserID == req.UserID {
			//我是发起方
			info = types.FriendInfoResponse{
				UserID:   int64(friend.RecvUserID),
				Nickname: friend.RecvUserModel.Nickname,
				Avatar:   friend.RecvUserModel.Avatar,
				Abstract: friend.RecvUserModel.Abstract,
				Notice:   friend.SendUserNotice,
			}
		}
		if friend.RecvUserID == req.UserID {
			//我是接收方
			info = types.FriendInfoResponse{

				UserID:   int64(friend.SendUserID),
				Nickname: friend.SendUserModel.Nickname,
				Avatar:   friend.SendUserModel.Avatar,
				Abstract: friend.SendUserModel.Abstract,
				Notice:   friend.RecvUserNotice,
			}
		}
		list = append(list, info)
	}
	return &types.FriendListResponse{
		Count: int(count),
		List:  list,
	}, nil

}
