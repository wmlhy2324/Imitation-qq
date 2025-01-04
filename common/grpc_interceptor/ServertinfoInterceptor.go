package grpc_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerInfoInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	clientIP := metadata.ValueFromIncomingContext(ctx, "clientIP")
	userID := metadata.ValueFromIncomingContext(ctx, "userID")
	if len(clientIP) > 0 {
		ctx = context.WithValue(ctx, "clientIP", clientIP)
	}
	if len(userID) > 0 {
		ctx = context.WithValue(ctx, "userID", userID)
	}
	return handler(ctx, req), nil
}
