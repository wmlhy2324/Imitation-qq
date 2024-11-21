package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64
	WhiteList []string
	BlackList []string
	UploadDir string //上传的文件目录
	UserRpc   zrpc.RpcClientConf
	Mysql     struct {
		DataSource string
	}
}
