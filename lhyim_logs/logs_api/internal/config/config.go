package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	KqConsumerConf kq.KqConf
	Etcd           string
	Mysql          struct {
		DataSource string
	}
	UserRpc zrpc.RpcClientConf
}
