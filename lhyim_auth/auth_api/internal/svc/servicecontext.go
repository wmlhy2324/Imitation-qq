package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/common/grpc_interceptor"
	"lhyim_server/common/log_stash"
	"lhyim_server/core"
	"lhyim_server/lhyim_auth/auth_api/internal/config"

	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	Redis          *redis.Client
	UserRpc        user_rpc.UsersClient
	KqPusherClient *kq.Pusher
	ActionLogs     *log_stash.Pusher
	RuntimeLogs    *log_stash.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	MysqlDB := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	kqClient := kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic)
	return &ServiceContext{
		Config:         c,
		DB:             MysqlDB,
		Redis:          client,
		UserRpc:        users.NewUsers(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(grpc_interceptor.ClientInfoInterceptor))),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		ActionLogs:     log_stash.NewActionPusher(kqClient, c.Name),
		RuntimeLogs:    log_stash.NewRuntimePusher(kqClient, c.Name),
	}
}
