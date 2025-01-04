package logic

import (
	"context"
	"encoding/json"
	"errors"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"
	"lhyim_server/lhyim_user/user_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type User_infoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUser_infoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *User_infoLogic {
	return &User_infoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *User_infoLogic) User_info(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	res, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoRequest{
		UserId: uint32(req.UserID),
	})
	if err != nil {
		return nil, err
	}
	var user user_models.UserModel
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("数据错误")
	}
	resp = &types.UserInfoResponse{
		UserID:        int64(user.ID),
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Abstract:      user.Abstract,
		ReCallMessage: user.UserConfModel.RecallMessage,
		FriendOnline:  user.UserConfModel.FriendOnline,
		Sound:         user.UserConfModel.Sound,
		SecureLink:    user.UserConfModel.SecureLink,
		SavePwd:       user.UserConfModel.SavePwd,
		SearchUser:    user.UserConfModel.SearchUser,
		Verification:  user.UserConfModel.Verification,
	}
	if user.UserConfModel.VerificationQuestion != nil {
		resp.VerificationQuestion = types.VerificationQuestion{
			Problem1: user.UserConfModel.VerificationQuestion.Problem1,
			Problem2: user.UserConfModel.VerificationQuestion.Problem2,
			Problem3: user.UserConfModel.VerificationQuestion.Problem3,
			Answer1:  user.UserConfModel.VerificationQuestion.Answer1,
			Answer2:  user.UserConfModel.VerificationQuestion.Answer2,
			Answer3:  user.UserConfModel.VerificationQuestion.Answer3,
		}
	}
	return resp, nil

}
