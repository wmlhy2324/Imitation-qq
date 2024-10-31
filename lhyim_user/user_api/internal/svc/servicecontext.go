package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/core"
	"lhyim_server/lhyim_user/user_api/internal/config"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UsersClient
	DB      *gorm.DB
	Redis   *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	MysqlDB := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		DB:      MysqlDB,
	}
}