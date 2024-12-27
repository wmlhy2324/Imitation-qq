package logic

import (
	"context"
	"fmt"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_group/group_models"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSessionLogic {
	return &GroupSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type SessionData struct {
	GroupID       uint   `gorm:"column:group_id"`
	NewMsgDate    string `gorm:"column:newMsgDate"`
	NewMsgPreview string `gorm:"column:newMsgPreview"`
}

func (l *GroupSessionLogic) GroupSession(req *types.GroupSessionRequest) (resp *types.GroupSessionResponse, err error) {
	//先查我有哪些群
	var userGroupIDList []uint
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).Where("user_id = ?", req.UserID).Select("group_id").Scan(&userGroupIDList)
	column := fmt.Sprintf("if((select %d from top_user_models where user_id = %d and (top_user_id = sU or top_user_id = rU) limit 1) ,1,0) as isTop", req.UserID, req.UserID)
	sessionList, count, _ := list_query.ListQuery(l.svcCtx.DB, SessionData{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "newMsgDate desc",
		},
		Table: func() (string, any) {
			return "(?) as u", l.svcCtx.DB.Model(&group_models.GroupMsgModel{}).
				Select("group_id",
					"max(created_at) as newMsgDate",
					"max(created_at) as maxDate",
					fmt.Sprintf("(select msg_preview from chat_models where ((send_user_id = sU and recv_user_id = rU) or (send_user_id = rU and recv_user_id = sU)) and id not in (select chat_id from user_chat_delete_models where user_id = %d) order by created_at desc limit 1) as maxPreview", req.UserID),
					column).
				Where("group_id in (?)", userGroupIDList).
				Group("group_id")
		},
	})
	var groupIDList []uint
	for _, data := range sessionList {
		groupIDList = append(groupIDList, data.GroupID)
	}
	var groupListModel []group_models.GroupModel
	l.svcCtx.DB.Find(&groupListModel, groupIDList)
	var groupMap = map[uint]group_models.GroupModel{}
	for _, model := range groupListModel {
		groupMap[model.ID] = model
	}
	resp = new(types.GroupSessionResponse)
	for _, data := range sessionList {
		resp.List = append(resp.List, types.GroupSessionList{
			GroupID:       data.GroupID,
			Title:         groupMap[data.GroupID].Title,
			Avatar:        groupMap[data.GroupID].Avatar,
			NewMsgDate:    data.NewMsgDate,
			NewMsgPreview: data.NewMsgPreview,
		})
	}
	resp.Count = int(count)
	return
}
