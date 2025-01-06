package main

import (
	"flag"
	"fmt"
	"lhyim_server/common/etcd"
	"lhyim_server/common/middleware"

	"lhyim_server/lhyim_auth/auth_api/internal/config"
	"lhyim_server/lhyim_auth/auth_api/internal/handler"
	"lhyim_server/lhyim_auth/auth_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/auth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.Use(middleware.LogActionMiddleware(ctx.ActionLogs))
	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
