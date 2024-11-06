package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_chat/chat_models"

	"lhyim_server/lhyim_chat/chat_rpc/internal/svc"
	"lhyim_server/lhyim_chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserChatLogic {
	return &UserChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserChatLogic) UserChat(in *chat_rpc.UserChatRequest) (*chat_rpc.UserChatResponse, error) {
	// todo: add your logic here and delete this line
	var msg ctype.Msg
	err := json.Unmarshal(in.Msg, &msg)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	var systemMsg *ctype.SystemMsg
	if in.SystemMsg != nil {
		err = json.Unmarshal(in.SystemMsg, systemMsg)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}
	fmt.Println("msg = ", msg)
	chat := chat_models.ChatModel{
		SendUserID: uint(in.SendUserId),
		RecvUserID: uint(in.RecvUserId),
		MsgType:    msg.Type,
		Msg:        msg,
		SystemMsg:  systemMsg,
	}
	chat.MsgPreview = chat.MsgPreviewMethod()
	fmt.Println(chat)
	err = l.svcCtx.DB.Create(&chat).Error
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	return &chat_rpc.UserChatResponse{}, nil
}
