package logic

import (
	"context"
	"fmt"
	"lhyim_server/lhyim_chat/chat_models"

	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatDeleteLogic {
	return &ChatDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatDeleteLogic) ChatDelete(req *types.ChatDeleteRequest) (resp *types.ChatDeleteResponse, err error) {

	var chatList []chat_models.ChatModel
	l.svcCtx.DB.Find(&chatList, req.IdList)

	var deleteChatList []chat_models.UserChatDeleteModel
	l.svcCtx.DB.Debug().Where("chat_id in (?) ", req.IdList).Find(&deleteChatList)
	chatDeleteMap := map[uint]struct{}{}
	for _, model := range deleteChatList {
		chatDeleteMap[model.ChatID] = struct{}{}
	}
	var userDeleteChatList []chat_models.UserChatDeleteModel
	if len(chatList) > 0 {
		for _, model := range chatList {
			//不是自己的聊天记录
			if !(model.SendUserID == req.UserID || model.RecvUserID == req.UserID) {
				continue
			}
			//已经删除过的聊天记录
			fmt.Println(model.ID)
			_, ok := chatDeleteMap[model.ID]
			fmt.Println(ok)
			if ok {
				continue
			}
			userDeleteChatList = append(userDeleteChatList, chat_models.UserChatDeleteModel{
				UserID: req.UserID,
				ChatID: model.ID,
			})
		}
	}
	if len(userDeleteChatList) > 0 {
		l.svcCtx.DB.Create(&userDeleteChatList)
	}
	logx.Infof("已删除聊天记录 %d 条", len(userDeleteChatList))
	return
}
