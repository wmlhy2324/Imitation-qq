package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/lhyim_auth/auth_models"
	"lhyim_server/utils/jwts"
	"lhyim_server/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/lhyim_auth/auth_api/internal/svc"
	"lhyim_server/lhyim_auth/auth_api/internal/types"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	fmt.Println(l.ctx.Value("clientIP"), l.ctx.Value("User-ID"))
	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "id = ?", req.Username).Error
	if err != nil {
		err = errors.New("用户不存在")
		return
	}
	if !pwd.CheckHash(user.Pwd, req.Password) {
		err = errors.New("用户名错误或者密码错误")
		return
	}
	//判断用户的登录来源，第三方登录的不能用密码进行登录
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserId:   int64(user.ID),
		Nickname: user.Nickname,
		Role:     user.Role,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Error(err)
		err = errors.New("服务器内部错误,生成token失败")
		return
	}
	err = l.svcCtx.KqPusherClient.Push(l.ctx, fmt.Sprintf("%s用户登录成功", user.Nickname))
	if err != nil {
		fmt.Println(err)
	}

	return &types.LoginResponse{Token: token}, nil
}
