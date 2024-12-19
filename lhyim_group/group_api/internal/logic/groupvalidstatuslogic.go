package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidStatusLogic {
	return &GroupValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidStatusLogic) GroupValidStatus(req *types.GroupValidStatusRequest) (resp *types.GroupValidStatusResponse, err error) {

	var groupValidModel group_models.GroupVerifyModel
	err = l.svcCtx.DB.Take(&groupValidModel, req.ValidID).Error
	if err != nil {
		return nil, errors.New("不存在的验证记录")
	}
	if groupValidModel.Status != 0 {
		return nil, errors.New("已经处理过该认证请求")
	}
	//判断我有没有权限处理请求
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&groupMember, "group_id = ? and user_id = ?", groupValidModel.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("群不存在或不是群成员")
	}
	if !(groupMember.Role == 1 || groupMember.Role == 2) {
		return nil, errors.New("无权限")
	}
	switch req.Status {
	case 0: //未操作
		break
	case 1: //1已同意
	//把用户加群里
	case 2: //2已拒绝
	case 3: //3已忽略

	}
	return
}
