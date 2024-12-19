package logic

import (
	"context"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidListLogic {
	return &GroupValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidListLogic) GroupValidList(req *types.GroupValidListRequest) (resp *types.GroupValidListResponse, err error) {
	//群验证列表,自己才是管理员或者群主才行
	var groupIDList []uint //
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).Where("user_id = ? and (role = 1 or role = 2)", req.UserID).Select("group_id").Scan(&groupIDList)
	//先去查自己管理了哪些群，去找这些群的验证表
	groups, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"GroupModel"},
		Where:   l.svcCtx.DB.Where("group_id in ?", groupIDList),
	})
	var userIDList []uint32
	for _, group := range groups {
		userIDList = append(userIDList, uint32(group.UserID))
	}
	userList, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	resp = new(types.GroupValidListResponse)
	resp.Count = int(count)
	for _, group := range groups {
		info := types.GroupValidiInfo{
			ID:               group.ID,
			GroupID:          group.GroupID,
			UserID:           group.UserID,
			Status:           group.Status,
			AddtionalMessage: group.AdditionalMessage,
			Title:            group.GroupModel.Title,
			CreateAt:         group.CreatedAt.String(),
			Type:             group.Type,
		}
		resp.List = append(resp.List)
		if group.VerificationQuestion != nil {
			info.VerificationQuestion = &types.VerificationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
				Answer1:  group.VerificationQuestion.Answer1,
				Answer2:  group.VerificationQuestion.Answer2,
				Answer3:  group.VerificationQuestion.Answer3,
			}
		}
		if err1 == nil {
			info.UserNickname = userList.UserInfo[uint32(group.UserID)].NickName
			info.UserAvatar = userList.UserInfo[uint32(group.UserID)].Avatar
		}
		resp.List = append(resp.List, info)
	}
	return
}
