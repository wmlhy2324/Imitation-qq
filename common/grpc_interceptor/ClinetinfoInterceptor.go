package grpc_interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var clientIP, userID string
	cl := ctx.Value("clientIP")
	if cl != nil {
		clientIP = cl.(string)
	}
	ui := ctx.Value("userID")
	if ui != nil {
		userID = ui.(string)
	}
	md := metadata.New(map[string]string{
		"clientIP": clientIP,
		"userID":   userID,
	})

	// Create a new context with the outgoing metadata
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}
	return nil
}
