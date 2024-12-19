package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRoleUpdateLogic {
	return &GroupMemberRoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRoleUpdateLogic) GroupMemberRoleUpdate(req *types.GroupMemberRoleUpdateRequest) (resp *types.GroupMemberRoleUpdateResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("你不是群成员")
	}
	if member.Role != 1 {
		return nil, errors.New("错误调用")
	}
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("该用户不是群成员")
	}
	if !(req.Role == 2 || req.Role == 3) {
		return nil, errors.New("用户角色错误")
	}
	if member1.Role == req.Role {
		return nil, nil
	}
	l.svcCtx.DB.Model(&member1).Update("role", req.Role)
	return
}
