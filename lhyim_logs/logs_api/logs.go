package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"lhyim_server/common/middleware"
	"lhyim_server/lhyim_logs/logs_api/internal/mqs"

	"lhyim_server/lhyim_logs/logs_api/internal/config"
	"lhyim_server/lhyim_logs/logs_api/internal/handler"
	"lhyim_server/lhyim_logs/logs_api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/logs.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	server.Use(middleware.LogMiddleware)
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	for _, mq := range mqs.Consumers(c, context.Background(), ctx) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()

	server.Start()
}
