package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberNicknameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberNicknameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberNicknameUpdateLogic {
	return &GroupMemberNicknameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberNicknameUpdateLogic) GroupMemberNicknameUpdate(req *types.GroupMemberNicknameUpdateRequest) (resp *types.GroupMemberNicknameUpdateResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("你不是群成员")
	}
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("该用户不是群成员")
	}
	//自己修改自己的
	if req.UserID == req.MemberID {
		l.svcCtx.DB.Model(&member).Updates(map[string]interface{}{
			"member_nickname": req.Nickname,
		})
		return
	}
	if !(member.Role == 1 && (member1.Role == 2 || member1.Role == 3) || (member.Role == 2 && member1.Role == 3)) {
		return nil, errors.New("用户角色错误")
	}
	l.svcCtx.DB.Model(&member1).Updates(map[string]interface{}{
		"member_nickname": req.Nickname,
	})
	return
}
