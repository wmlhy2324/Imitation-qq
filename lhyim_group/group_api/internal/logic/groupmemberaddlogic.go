package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberAddLogic {
	return &GroupMemberAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberAddLogic) GroupMemberAdd(req *types.GroupMemberAddRequest) (resp *types.GroupMemberAddResponse, err error) {
	//群成员邀请好友
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("本群成员才可以邀请用户")
	}
	if member.Role == 3 {
		if !member.GroupModel.IsInvite {
			return nil, errors.New("管理员未开放好友入群功能")
		}
	}
	//查哪些用户已经进群了
	var memberList []group_models.GroupMemberModel

	l.svcCtx.DB.Find(&memberList, "group_id = ? and user_id in ?", req.ID, req.MemberIDList)
	if len(memberList) > 0 {
		return nil, errors.New("已经有用户在群里了")
	}
	for _, memberID := range req.MemberIDList {
		memberList = append(memberList, group_models.GroupMemberModel{
			GroupID: req.ID,
			UserID:  memberID,
			Role:    3,
		})
	}
	err = l.svcCtx.DB.Create(&memberList).Error
	if err != nil {
		logx.Error(err)
	}
	return
}
