package logic

import (
	"context"
	"errors"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	pwd2 "lhyim_server/utils/pwd"

	"lhyim_server/lhyim_auth/auth_api/internal/svc"
	"lhyim_server/lhyim_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	pwd := pwd2.HashPwd(req.Pwd)
	_, err = l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateRequest{
		NickName:       req.Username,
		Password:       pwd,
		Role:           1,
		Avatar:         "",
		OpenId:         "001",
		RegisterSource: "本地注册",
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("注册失败")
	}
	return
}
