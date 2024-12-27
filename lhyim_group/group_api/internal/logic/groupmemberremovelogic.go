package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRemoveLogic {
	return &GroupMemberRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRemoveLogic) GroupMemberRemove(req *types.GroupMemberRemoveRequest) (resp *types.GroupMemberRemoveResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("你不是群成员")
	}
	//用户自己退群，群主不能退群，群主只能解散群
	if req.UserID == req.MemberID {
		if member.Role == 1 {
			return nil, errors.New("群主不能退群，只能解散群")
		}
		//把member中与这个用户的记录删掉
		l.svcCtx.DB.Delete(&member)
		//给群验证表里面加记录
		l.svcCtx.DB.Create(&group_models.GroupVerifyModel{
			GroupID: member.GroupID,
			UserID:  req.UserID,
			Type:    2,
		})
	}
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("非法调用")
	}
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("该用户不是群成员")
	}
	//群出可以踢管理员和用户//管理员可以踢用户
	if !(member.Role == 1 && (member1.Role == 2 || member1.Role == 3) || (member.Role == 2 && member1.Role == 3)) {
		return nil, errors.New("角色错误")
	}
	err = l.svcCtx.DB.Delete(&member1).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群成员移除失败")
	}

	return
}
