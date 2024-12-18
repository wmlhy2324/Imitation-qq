package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupRemoveLogic {
	return &GroupRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupRemoveLogic) GroupRemove(req *types.GroupRemoveRequest) (resp *types.GroupRemoveResponse, err error) {

	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&groupMember, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("群不存在或者用户不是群成员")
	}
	if groupMember.Role != 1 {
		return nil, errors.New("只有群主能解散")
	}
	//这个群关联的群消息要删掉
	var msgList []group_models.GroupMsgModel
	l.svcCtx.DB.Find(&msgList, "group_id = ?", req.ID).Delete(&msgList)
	//群成员删掉
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID).Delete(&memberList)
	//群验证消息
	var vList []group_models.GroupVerifyModel
	l.svcCtx.DB.Find(&vList, "group_id = ?", req.ID).Delete(&vList)
	//群删掉
	var group group_models.GroupModel
	l.svcCtx.DB.Take(&group, req.ID).Delete(&group)

	logx.Infof("删除群:%s", group.Title)
	return
}
