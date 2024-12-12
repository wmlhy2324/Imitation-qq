package svc

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"lhyim_server/core"
	"lhyim_server/lhyim_group/group_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	MysqlDB := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.DB)
	return &ServiceContext{
		Config: c,
		DB:     MysqlDB,
		Redis:  client,
	}
}
