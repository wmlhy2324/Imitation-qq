package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserValidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserValidLogic {
	return &UserValidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserValidLogic) UserValid(req *types.UserValidRequest) (resp *types.UserValidResponse, err error) {
	// todo: add your logic here and delete this line
	//如果是好久就不用加了
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, req.FriendID, req.UserID) {
		return nil, errors.New("你们已经是好友了")
	}
	var userConf user_models.UserConfModel
	err = l.svcCtx.DB.Take(&userConf, "user_id=?", req.FriendID).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	resp = new(types.UserValidResponse)
	resp.Verification = userConf.Verification
	switch userConf.Verification {
	case 0: //不允许添加
	case 1: //允许任何人
	case 2: //需要验证问题
	case 3, 4: //需要回答问题
		if userConf.VerificationQuestion != nil {
			resp.VerificationQuestion = types.VerificationQuestion{
				Problem1: userConf.VerificationQuestion.Problem1,
				Problem2: userConf.VerificationQuestion.Problem2,
				Problem3: userConf.VerificationQuestion.Problem3,
			}
		}
	default:

	}

	return
}
