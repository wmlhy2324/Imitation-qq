// Code generated by goctl. DO NOT EDIT.
// Source: user_rpc.proto

package server

import (
	"context"

	"lhyim_server/lhyim_user/user_rpc/internal/logic"
	"lhyim_server/lhyim_user/user_rpc/internal/svc"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
)

type UsersServer struct {
	svcCtx *svc.ServiceContext
	user_rpc.UnimplementedUsersServer
}

func NewUsersServer(svcCtx *svc.ServiceContext) *UsersServer {
	return &UsersServer{
		svcCtx: svcCtx,
	}
}

func (s *UsersServer) UserCreate(ctx context.Context, in *user_rpc.UserCreateRequest) (*user_rpc.UserCreateResponse, error) {
	l := logic.NewUserCreateLogic(ctx, s.svcCtx)
	return l.UserCreate(in)
}

func (s *UsersServer) UserInfo(ctx context.Context, in *user_rpc.UserInfoRequest) (*user_rpc.UserInfoResponse, error) {
	l := logic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}

func (s *UsersServer) UserListInfo(ctx context.Context, in *user_rpc.UserListInfoRequest) (*user_rpc.UserListInfoResponse, error) {
	l := logic.NewUserListInfoLogic(ctx, s.svcCtx)
	return l.UserListInfo(in)
}

func (s *UsersServer) IsFriend(ctx context.Context, in *user_rpc.IsFriendRequest) (*user_rpc.IsFriendResponse, error) {
	l := logic.NewIsFriendLogic(ctx, s.svcCtx)
	return l.IsFriend(in)
}

func (s *UsersServer) FriendList(ctx context.Context, in *user_rpc.FriendListRequest) (*user_rpc.FriendListResponse, error) {
	l := logic.NewFriendListLogic(ctx, s.svcCtx)
	return l.FriendList(in)
}

func (s *UsersServer) UserBaseInfo(ctx context.Context, in *user_rpc.UserBaseInfoRequest) (*user_rpc.UserBaseInfoResponse, error) {
	l := logic.NewUserBaseInfoLogic(ctx, s.svcCtx)
	return l.UserBaseInfo(in)
}

func (s *UsersServer) UserOnlineList(ctx context.Context, in *user_rpc.UserOnlineRequest) (*user_rpc.UserOnlineResponse, error) {
	l := logic.NewUserOnlineListLogic(ctx, s.svcCtx)
	return l.UserOnlineList(in)
}
