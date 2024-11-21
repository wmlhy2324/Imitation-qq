package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"lhyim_server/core"
	"lhyim_server/lhyim_chat/chat_api/internal/config"
	"lhyim_server/lhyim_file/file_rpc/files"
	"lhyim_server/lhyim_file/file_rpc/types/file_rpc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/lhyim_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	UserRpc user_rpc.UsersClient
	Redis   *redis.Client
	FileRpc file_rpc.FilesClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDB := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDB,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		Redis:   client,
		FileRpc: files.NewFiles(zrpc.MustNewClient(c.FileRpc)),
	}
}
