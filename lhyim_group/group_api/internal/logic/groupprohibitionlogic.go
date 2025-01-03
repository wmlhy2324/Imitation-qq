package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/lhyim_group/group_models"
	"time"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupProhibitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupProhibitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupProhibitionLogic {
	return &GroupProhibitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupProhibitionLogic) GroupProhibition(req *types.GroupProhibitionRequest) (resp *types.GroupProhibitionResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "user_id = ? and group_id = ?", req.UserID, req.GroupID).Error
	if err != nil {
		return nil, errors.New("群不存在或者你不是群成员")
	}
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("当前用户角色错误")
	}

	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "user_id = ? and group_id = ?", req.MemberID, req.GroupID).Error
	if err != nil {
		return nil, errors.New("群不存在或者该用户不是群成员")
	}
	if !(member.Role == 1 && ((member1.Role == 2) || (member1.Role == 3)) || (member.Role == 2 && member1.Role == 3)) {

		return nil, errors.New("角色权限错误")
	}
	l.svcCtx.DB.Model(&member1).Update("prohibition_time", req.ProhibitionTime)

	//利用redis的过期时间去做禁言时间
	key := fmt.Sprintf("prohibition_%d", member1.ID)
	if req.ProhibitionTime != nil {
		//给redis设置一个key，过期时间是

		l.svcCtx.Redis.Set(key, "time", time.Duration(*req.ProhibitionTime)*time.Second)
	} else {
		l.svcCtx.Redis.Del(key)
	}
	return
}
