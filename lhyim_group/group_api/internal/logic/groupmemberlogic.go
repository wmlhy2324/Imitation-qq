package logic

import (
	"context"
	"fmt"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberLogic {
	return &GroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Data struct {
	GroupID        uint   `gorm:"column:group_id"`
	UserID         uint   `gorm:"column:user_id"`
	Role           string `gorm:"column:role"`
	CreatedAt      string `gorm:"column:created_at"`
	MemberNickname string `gorm:"column:member_nickname"`
	NewMsgDate     string `gorm:"column:new_msg_date"`
}

func (l *GroupMemberLogic) GroupMember(req *types.GroupMemberRequest) (resp *types.GroupMemberResponse, err error) {
	// todo: add your logic here and delete this line

	column := fmt.Sprintf("(select group_msg_models.created_at from group_msg_models where group_member_models.group_id = %d and group_msg_models.send_user_id = user_id ) as new_msg_date", req.ID)
	memberList, count, _ := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		},
		Table: func() (string, any) {
			return "(?) as u", l.svcCtx.DB.Model(&group_models.GroupMemberModel{GroupID: req.ID}).
				Select(
					"group_id",
					"user_id",
					"role",
					"created_at",
					"member_nickname",
					column,
				)
		},
	})
	var userIDList []uint32
	for _, data := range memberList {
		userIDList = append(userIDList, uint32(data.UserID))
	}
	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{UserIdList: userIDList})
	//服务降级
	if err != nil {
		logx.Error(err)
		//如果这里不初始化的话，后面取值会报错空指针，导致整个服务挂掉
		userListResponse.UserInfo = map[uint32]*user_rpc.UserInfo{}
	}
	var userOnlineMap = map[uint32]bool{}

	userOnlineResponse, err := l.svcCtx.UserRpc.UserOnlineList(l.ctx, &user_rpc.UserOnlineRequest{})
	if err == nil {
		for _, v := range userOnlineResponse.UserIdList {
			userOnlineMap[v] = true
		}
	}
	return
}
