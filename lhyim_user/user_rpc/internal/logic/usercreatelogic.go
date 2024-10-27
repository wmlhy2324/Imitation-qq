package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_rpc/internal/svc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *user_rpc.UserCreateRequest) (*user_rpc.UserCreateResponse, error) {
	// todo: add your logic here and delete this line
	var user user_models.UserModel
	err := l.svcCtx.DB.Take(&user, "open_id = ?", in.OpenId).Error
	if err == nil {
		return nil, errors.New("用户已存在")
	}
	user = user_models.UserModel{
		Nickname:       in.NickName,
		OpenID:         in.OpenId,
		Avatar:         in.Avatar,
		Role:           int8(2),
		RegisterSource: in.RegisterSource,
		Pwd:            in.Password,
	}
	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		return nil, errors.New("用户创建失败")
	}
	//创建用户配置
	l.svcCtx.DB.Create(&user_models.UserConfModel{
		UserID:        user.ID,
		RecallMessage: nil,
		FriendOnline:  false, //关闭好友上线提醒
		Sound:         true,
		SecureLink:    false,
		SavePwd:       false,
		SearchUser:    2,
		Verification:  2, //需要验证消息
		Online:        true,
	})
	return &user_rpc.UserCreateResponse{UserId: int32(user.ID)}, nil
}
