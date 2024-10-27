package logic

import (
	"context"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"lhyim_server/lhyim_chat/chat_models"
	"lhyim_server/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
type ChatHistory struct {
	ID        uint             `json:"id"`
	SendUser  UserInfo         `json:"sendUser"`
	RevUser   UserInfo         `json:"revUser"`
	IsMe      bool             `json:"isMe"` //哪条消息是我发的
	CreateAt  string           `json:"createAt"`
	Msg       ctype.Msg        `json:"msg"`
	SystemMsg *ctype.SystemMsg `json:"systemMsg"`
}
type ChatHistoryResponse struct {
	List  []ChatHistory `json:"list"`
	Count int64         `json:"count"`
}

func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryRequest) (resp *ChatHistoryResponse, err error) {

	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatModel{},
		list_query.Option{
			PageInfo: models.PageInfo{
				Page:  req.Page,
				Limit: req.Limit,
			},
			Where: l.svcCtx.DB.Where("send_user_id = ? or recv_user_id = ?", req.UserID, req.UserID),
		})
	var userIDList []uint
	for _, model := range chatList {
		userIDList = append(userIDList, model.SendUserID)
		userIDList = append(userIDList, model.RecvUserID)
	}
	//去重
	userIDList = utils.DeduplicatoionList(userIDList)
	//去调用户服务的rpc方法，获取用户信息
	var list = make([]ChatHistory, 0)
	for _, model := range chatList {
		list = append(list, ChatHistory{
			ID:        model.ID,
			CreateAt:  model.CreatedAt,
			Msg:       model.Msg,
			SystemMsg: model.SystemMsg,
		})
	}
	resp = &ChatHistoryResponse{
		Count: count,
		List:  list,
	}
	return
}
