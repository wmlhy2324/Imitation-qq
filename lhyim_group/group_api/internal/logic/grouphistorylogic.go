package logic

import (
	"context"
	"errors"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryLogic {
	return &GroupHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type HistoryResponse struct {
	UserID         uint      `json:"userID"`
	UserNickname   string    `json:"userNickname"`
	UserAvatar     string    `json:"userAvatar"`
	Msg            ctype.Msg `json:"msg"`
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	MsgType        int8      `json:"msgType"`
	IsMe           bool      `json:"isMe"`
	MemberNickname string    `json:"memberNickname"` //群用户备注
}
type HistoryListResponse struct {
	List  []HistoryResponse `json:"list"`
	Count int               `json:"count"`
}

func (l *GroupHistoryLogic) GroupHistory(req *types.GroupHistoryRequest) (resp *HistoryListResponse, err error) {
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}
	var msgIDList = []uint{0}
	l.svcCtx.DB.Model(&group_models.GroupUserMsgDeleteModel{}).Select("msg_id").
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList)
	groupMsgs, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupMsgModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where:   l.svcCtx.DB.Where("id not in ?", msgIDList),
		Preload: []string{"MemberModel"},
	})

	var userIDList []uint32
	for _, model := range groupMsgs {
		userIDList = append(userIDList, uint32(model.SendUserID))
	}
	userIDList = utils.DeduplicatoionList(userIDList)
	userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	var list = make([]HistoryResponse, 0)
	for _, model := range groupMsgs {
		info := HistoryResponse{
			UserID:    model.SendUserID,
			Msg:       model.Msg,
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			MsgType:   int8(model.MsgType),
		}
		if model.MemberModel != nil {
			info.MemberNickname = model.MemberModel.MemberNickname
		}
		if err1 == nil {
			info.UserNickname = userListResponse.UserInfo[uint32(model.SendUserID)].NickName
			info.UserAvatar = userListResponse.UserInfo[uint32(model.SendUserID)].Avatar
		}
		if req.UserID == info.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}
	resp = new(HistoryListResponse)
	resp.List = list
	resp.Count = int(count)
	return
}
