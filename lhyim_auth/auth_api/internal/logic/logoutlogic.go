package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/utils/jwts"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/lhyim_auth/auth_api/internal/svc"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(token string) (resp string, err error) {
	// todo: add your logic here and delete this line
	if token == "" {
		err = errors.New("token不能为空")
		return
	}
	payload, err := jwts.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("token解析失败")
		return
	}
	now := time.Now()
	expiration := payload.ExpiresAt.Time.Sub(now)
	key := fmt.Sprintf("logout_%s", token)
	l.svcCtx.Redis.SetNX(key, "", expiration)
	resp = "注销登录成功"
	return resp, nil
}
