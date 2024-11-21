package logic

import (
	"context"
	"encoding/json"
	"errors"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_chat/chat_rpc/chat"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidStatusLogic {
	return &ValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidStatusLogic) ValidStatus(req *types.FriendValidStatusRequest) (resp *types.FriendValidStatusResponse, err error) {
	var friendverify user_models.FriendVerifyModel
	err = l.svcCtx.DB.Take(&friendverify, "id = ? and recv_user_id = ?", req.VerifyID, req.UserID).Error
	if err != nil {
		return nil, errors.New("验证记录不存在")
	}
	if friendverify.RecvStatus != 0 {
		return nil, errors.New("不可更改状态")
	}

	switch req.Status {
	case 1: //同意
		friendverify.RecvStatus = 1
		//往好友表里面加
		l.svcCtx.DB.Create(&user_models.FriendModel{
			SendUserID: friendverify.SendUserID,
			RecvUserID: friendverify.RecvUserID,
		})
		//给对方发消息
		msg := ctype.Msg{
			Type: ctype.TextMsgType,
			TextMsg: &ctype.TextMsg{
				Content: "我们已经是好友了,快来一起聊天吧",
			},
		}
		byteData, _ := json.Marshal(msg)
		_, err := l.svcCtx.ChatRpc.UserChat(context.Background(), &chat.UserChatRequest{
			SendUserId: uint32(friendverify.SendUserID),
			RecvUserId: uint32(friendverify.RecvUserID),
			Msg:        byteData,
			SystemMsg:  nil,
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	case 2: //拒绝
		friendverify.RecvStatus = 2
	case 3: //忽略
		friendverify.RecvStatus = 3
	case 4: //删除
		//一条验证记录是两个人看的
		l.svcCtx.DB.Model(&friendverify).Delete(&friendverify)
		return nil, nil
	}
	l.svcCtx.DB.Save(&friendverify)
	return
}
