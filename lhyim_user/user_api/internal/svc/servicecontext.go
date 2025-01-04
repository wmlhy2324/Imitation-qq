package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/common/grpc_interceptor"
	"lhyim_server/core"
	"lhyim_server/lhyim_chat/chat_rpc/chat"
	"lhyim_server/lhyim_chat/chat_rpc/types/chat_rpc"
	"lhyim_server/lhyim_user/user_api/internal/config"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UsersClient
	DB      *gorm.DB
	Redis   *redis.Client
	ChatRpc chat_rpc.ChatClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	MysqlDB := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	return &ServiceContext{
		Config:  c,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
		DB:      MysqlDB,
		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
		Redis:   client,
	}
}
