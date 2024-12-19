package logic

import (
	"context"
	"fmt"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils/set"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSearchLogic {
	return &GroupSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupSearchLogic) GroupSearch(req *types.GroupSearchRequest) (resp *types.GroupSearchListResponse, err error) {
	//insearch为false表示不能被搜索
	groups, _, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"MemberList"},
		Where:   l.svcCtx.DB.Where("is_search = 1 and (id = ? or title like ?)", req.Key, fmt.Sprintf("%%%s%%", req.Key)),
	})
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineRequest{})
	var userOnlineIDList []uint
	if err == nil {
		//降级
		for _, u := range userOnlineResponse.UserIdList {
			userOnlineIDList = append(userOnlineIDList, uint(u))
		}
	}
	resp = new(types.GroupSearchListResponse)
	for _, group := range groups {

		var groupMemberIdList []uint
		var isInGroup bool
		for _, member := range group.MemberList {
			groupMemberIdList = append(groupMemberIdList, member.UserID)
			if member.UserID == req.UserID {
				isInGroup = true
			}
		}
		//算在线总数

		resp.List = append(resp.List, types.GroupSearch{
			GroupID:         group.ID,
			Title:           group.Title,
			Abstract:        group.Abstract,
			Avatar:          group.Avatar,
			UserCount:       len(group.MemberList),
			UserOnlineCount: len(set.Intersect(groupMemberIdList, userOnlineIDList)), //在线用户总数
			IsInGroup:       isInGroup,                                               //我是否在群里面
		})
	}
	return
}
