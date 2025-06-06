// Code generated by goctl. DO NOT EDIT.
// Source: file_rpc.proto

package files

import (
	"context"

	"lhyim_server/lhyim_file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	FileInfoRequest  = file_rpc.FileInfoRequest
	FileInfoResponse = file_rpc.FileInfoResponse

	Files interface {
		FileInfo(ctx context.Context, in *FileInfoRequest, opts ...grpc.CallOption) (*FileInfoResponse, error)
	}

	defaultFiles struct {
		cli zrpc.Client
	}
)

func NewFiles(cli zrpc.Client) Files {
	return &defaultFiles{
		cli: cli,
	}
}

func (m *defaultFiles) FileInfo(ctx context.Context, in *FileInfoRequest, opts ...grpc.CallOption) (*FileInfoResponse, error) {
	client := file_rpc.NewFilesClient(m.cli.Conn())
	return client.FileInfo(ctx, in, opts...)
}
