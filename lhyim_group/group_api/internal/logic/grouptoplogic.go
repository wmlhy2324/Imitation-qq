package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupTopLogic {
	return &GroupTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupTopLogic) GroupTop(req *types.GroupTopRequest) (resp *types.GroupTopResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("你还不是群成员")
	}
	var userTop group_models.GroupUserTopModel
	err1 := l.svcCtx.DB.Take(&userTop, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err1 != nil {
		if req.IsTop {
			l.svcCtx.DB.Create(&group_models.GroupUserTopModel{
				GroupID: req.GroupID,
				UserID:  req.UserID,
			})
		}
	} else {
		if !req.IsTop {
			l.svcCtx.DB.Delete(&userTop)
		}
	}
	return
}
