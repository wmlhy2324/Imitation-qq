package logic

import (
	"context"
	"errors"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"lhyim_server/lhyim_chat/chat_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
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
				Sort:  "created_at desc",
				Page:  req.Page,
				Limit: req.Limit,
			},
			Where: l.svcCtx.DB.Where("send_user_id = ? or recv_user_id = ?", req.UserID, req.UserID),
		})
	var userIDList []uint32
	for _, model := range chatList {
		userIDList = append(userIDList, uint32(model.SendUserID))
		userIDList = append(userIDList, uint32(model.RecvUserID))
	}
	//去重
	userIDList = utils.DeduplicatoionList(userIDList)
	//去调用户服务的rpc方法，获取用户信息
	response, err := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{UserIdList: userIDList})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}
	var list = make([]ChatHistory, 0)
	for _, model := range chatList {
		sendUser := UserInfo{
			ID:       model.SendUserID,
			Nickname: response.UserInfo[uint32(model.SendUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.SendUserID)].Avatar,
		}
		revUser := UserInfo{
			ID:       model.RecvUserID,
			Nickname: response.UserInfo[uint32(model.RecvUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.RecvUserID)].Avatar,
		}
		info := ChatHistory{
			ID:        model.ID,
			SendUser:  sendUser,
			RevUser:   revUser,
			CreateAt:  model.CreatedAt.String(),
			Msg:       model.Msg,
			SystemMsg: model.SystemMsg,
		}
		if info.SendUser.ID == req.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}
	resp = &ChatHistoryResponse{
		Count: count,
		List:  list,
	}
	return
}
