package logic

import (
	"context"

	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHanderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHanderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHanderLogic {
	return &ChatHanderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatHanderLogic) ChatHander(req *types.ChatRequest) (resp *types.ChatDeleteResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
