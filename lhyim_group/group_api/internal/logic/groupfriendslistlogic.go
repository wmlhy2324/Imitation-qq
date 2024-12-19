package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupfriendsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupfriendsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupfriendsListLogic {
	return &GroupfriendsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupfriendsListLogic) GroupfriendsList(req *types.GroupfriendsListRequest) (resp *types.GroupfriendsListResponse, err error) {
	//我的好友哪些在这个群
	friendResponse, err := l.svcCtx.UserRpc.FriendList(l.ctx, &user_rpc.FriendListRequest{
		UserId: uint32(req.UserID),
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("好友查询错误")
	}
	//需要去查我的好友列表

	//这个群的群成员列表，租成一个map
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID)
	var memberMap = map[uint]bool{}
	for _, member := range memberList {
		memberMap[member.UserID] = true

	}
	resp = new(types.GroupfriendsListResponse)
	for _, info := range friendResponse.FriendList {
		resp.List = append(resp.List, types.GroupfriendsList{
			UserId:    uint(info.UserId),
			Avatar:    info.Avatar,
			Nickname:  info.NickName,
			IsInGroup: memberMap[uint(info.UserId)],
		})
	}
	return
}
