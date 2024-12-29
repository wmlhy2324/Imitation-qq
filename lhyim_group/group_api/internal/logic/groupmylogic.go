package logic

import (
	"context"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMyLogic {
	return &GroupMyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMyLogic) GroupMy(req *types.GroupMyRequest) (resp *types.GroupMyListResponse, err error) {
	//查群id列表
	var groupIDList []uint
	query := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).Where("user_id = ?", req.UserID)
	if req.Mode == 1 {
		//我创建
		query.Where("role = ?", 1)
	}
	query.Select("group_id").Scan(&groupIDList)
	groupList, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{
			"MemberList",
		},
		Where: l.svcCtx.DB.Where("id in ?", groupIDList),
	})
	resp = new(types.GroupMyListResponse)
	for _, model := range groupList {
		var role int8
		for _, memberModel := range model.MemberList {
			if req.UserID == memberModel.UserID {
				role = memberModel.Role
			}
		}
		resp.List = append(resp.List, types.GroupMyResponse{
			GroupID:          model.ID,
			GroupAvatar:      model.Avatar,
			GroupTitle:       model.Title,
			GroupMemberCount: len(model.MemberList),
			Role:             role,
			Mole:             req.Mode,
		})
	}
	resp.Count = int(count)
	return
}
