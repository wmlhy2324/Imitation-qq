package svc

import (
	"gorm.io/gorm"
	"lhyim_server/core"
	"lhyim_server/lhyim_file/file_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	MysqlDB := core.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		DB:     MysqlDB,
	}
}
