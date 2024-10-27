package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/utils"
	"lhyim_server/utils/jwts"

	"lhyim_server/lhyim_auth/auth_api/internal/svc"
	"lhyim_server/lhyim_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationRequest) (resp *types.AuthenticationResponse, err error) {
	// todo: add your logic here and delete this line
	if utils.ListByRegex(l.svcCtx.Config.WriteList, req.ValidPath) {
		logx.Infof("请求地址：%s,在白名单中", req.ValidPath)
		return
	}
	if req.Token == "" {
		err = errors.New("认证失败,token不能为空")
		return
	}
	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("认证失败,token解析失败")
		return
	}
	_, err = l.svcCtx.Redis.Get(fmt.Sprintf("logout_%s", req.Token)).Result()
	if err == nil {
		err = errors.New("认证失败,token已注销")
		return
	}
	err = nil

	return &types.AuthenticationResponse{
		UserID: claims.UserId,
		Role:   int(claims.Role),
	}, nil
}
