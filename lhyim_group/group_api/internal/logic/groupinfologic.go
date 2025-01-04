package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils/set"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoLogic {
	return &GroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInfoLogic) GroupInfo(req *types.GroupInfoRequest) (resp *types.GroupInfoResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error

	if err != nil {
		logx.Error(err)
		return nil, errors.New("该用户不是群成员")
	}

	var groupModel group_models.GroupModel
	err = l.svcCtx.DB.Preload("MemberList").Take(&groupModel, req.ID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群不存在")
	}
	resp = &types.GroupInfoResponse{
		GroupID:         groupModel.ID,
		Title:           groupModel.Title,
		MemberCount:     len(groupModel.MemberList),
		Avatar:          groupModel.Avatar,
		Abstract:        groupModel.Abstract,
		Role:            member.Role,
		IsProhibition:   groupModel.IsProhibition,
		ProhibitionTime: member.GetProhibitionTime(l.svcCtx.Redis, l.svcCtx.DB),
	}
	//用户信息列表
	var userIDList []uint32
	var userAllIDList []uint32
	for _, member := range groupModel.MemberList {
		if member.Role == 1 || member.Role == 2 {
			userIDList = append(userIDList, uint32(member.UserID))
		}
		userAllIDList = append(userIDList, uint32(member.UserID))
	}
	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{UserIdList: userIDList})

	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户信息列表rpc错误")
	}
	//用户列表信息
	var Creator types.UserInfo
	var AdminList []types.UserInfo
	for _, model := range groupModel.MemberList {
		if model.Role == 3 {
			continue
		}
		userinfo := types.UserInfo{
			UserID:   model.UserID,
			Avatar:   userListResponse.UserInfo[uint32(model.UserID)].Avatar,
			Nickname: userListResponse.UserInfo[uint32(model.UserID)].NickName,
		}
		if model.Role == 1 {
			Creator = userinfo
		}
		if model.Role == 2 {
			AdminList = append(AdminList, userinfo)
		}
	}
	resp.Creator = Creator
	resp.AdminList = AdminList
	//算在线用户总个数
	//用户服务需要去写一个
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineRequest{})
	if err == nil {
		//算群成员与总在线人员取交集

		slice := set.Intersect(userOnlineResponse.UserIdList, userAllIDList)
		resp.MemberOnlineCount = len(slice)

	}
	return
}
