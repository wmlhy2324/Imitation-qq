package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/lhyim_auth/auth_api/internal/svc"
	"lhyim_server/lhyim_auth/auth_api/internal/types"
	"lhyim_server/lhyim_auth/auth_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils/jwts"
	"lhyim_server/utils/open_login"

	"github.com/zeromicro/go-zero/core/logx"
)

type Open_loginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpen_loginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Open_loginLogic {
	return &Open_loginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Open_loginLogic) Open_login(req *types.OpenLoginRequest) (resp *types.LoginResponse, err error) {
	type OpenIndo struct {
		Nickname string
		Avatar   string
		OpenID   string
	}
	var info OpenIndo
	switch req.Flag {
	case "qq":
		qqinfo, OpenError := open_login.NewQQLogin(req.Code, open_login.QQConfig{
			AppID:    l.svcCtx.Config.QQ.AppID,
			AppKey:   l.svcCtx.Config.QQ.AppKey,
			Redirect: l.svcCtx.Config.QQ.Redirect,
		})

		info = OpenIndo{
			Nickname: qqinfo.Nickname,
			Avatar:   qqinfo.Avatar,
			OpenID:   qqinfo.OpenID,
		}
		fmt.Println(info)
		err = OpenError
	default:
		err = errors.New("不支持的登录方式")
	}
	if err != nil {
		logx.Error(err)
		return nil, errors.New("登录失败")
	}
	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "open_id = ?", info.OpenID).Error
	if err != nil {
		//如果没有查到用户，说明是第一次登录，需要注册
		fmt.Println("注册逻辑")
		res, err := l.svcCtx.UserRpc.UserCreate(l.ctx, &user_rpc.UserCreateRequest{

			NickName:       info.Nickname,
			Password:       "",
			Role:           1,
			Avatar:         info.Avatar,
			OpenId:         info.OpenID,
			RegisterSource: "qq",
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("注册失败")
		}
		user.Model.ID = uint(res.UserId)
		user.Role = 2
		user.Nickname = info.Nickname
	}
	token, err := jwts.GenToken(jwts.JwtPayload{
		UserId:   int64(user.ID),
		Nickname: user.Nickname,
		Role:     user.Role,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Error(err)
		err = errors.New("服务器内部错误,生成token失败")
		return &types.LoginResponse{Token: token}, nil
	}
	return
}
