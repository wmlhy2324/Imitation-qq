package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/common/grpc_interceptor"
	"lhyim_server/core"
	"lhyim_server/lhyim_group/group_api/internal/config"
	"lhyim_server/lhyim_group/group_rpc/groups"
	"lhyim_server/lhyim_group/group_rpc/types/group_rpc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	UserRpc  user_rpc.UsersClient
	GroupRpc group_rpc.GroupsClient
	Redis    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDB := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	return &ServiceContext{
		Config:   c,
		DB:       mysqlDB,
		UserRpc:  users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
		Redis:    client,
		GroupRpc: groups.NewGroups(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
	}
}
