package main

import (
	"flag"
	"fmt"
	"lhyim_server/common/grpc_interceptor"

	"lhyim_server/lhyim_chat/chat_rpc/internal/config"
	"lhyim_server/lhyim_chat/chat_rpc/internal/server"
	"lhyim_server/lhyim_chat/chat_rpc/internal/svc"
	"lhyim_server/lhyim_chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/chatrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		chat_rpc.RegisterChatServer(grpcServer, server.NewChatServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	s.AddUnaryInterceptors(grpc_interceptor.ServerInfoInterceptor)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
