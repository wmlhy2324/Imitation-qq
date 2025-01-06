package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/common/grpc_interceptor"
	"lhyim_server/core"
	"lhyim_server/lhyim_logs/logs_api/internal/config"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	UserRpc user_rpc.UsersClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDB := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
		DB:      mysqlDB,
	}
}
