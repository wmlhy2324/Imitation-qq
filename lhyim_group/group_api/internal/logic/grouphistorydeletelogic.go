package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/utils/set"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryDeleteLogic {
	return &GroupHistoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupHistoryDeleteLogic) GroupHistoryDelete(req *types.GroupHistoryDeleteRequest) (resp *types.GroupHistoryDeleteResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}
	//删除的聊天记录
	var msgIDList []uint
	l.svcCtx.DB.Model(&group_models.GroupUserMsgDeleteModel{}).Select("msg_id").
		Where("group_id = ? and user_id = ?", req.GroupID, req.UserID).Scan(&msgIDList)

	addMsgIDList := set.Difference(req.MsgIDList, msgIDList)

	//用户传过来的消息id不一定存在
	var IsMgsExist []uint
	l.svcCtx.DB.Model(&group_models.GroupMsgModel{}).Where("id IN (?)", addMsgIDList).
		Select("id").Scan(&IsMgsExist)
	if len(IsMgsExist) != len(addMsgIDList) {
		return nil, errors.New("消息一致性异常")
	}

	var List []group_models.GroupUserMsgDeleteModel
	for _, msgID := range addMsgIDList {
		List = append(List, group_models.GroupUserMsgDeleteModel{
			UserID:  req.UserID,
			MsgID:   msgID,
			GroupID: req.GroupID,
		})
	}
	err = l.svcCtx.DB.Create(&List).Error
	if err != nil {
		return nil, errors.New("删除聊天记录失败")
	}
	return
}
