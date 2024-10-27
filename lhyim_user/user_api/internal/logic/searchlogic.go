package logic

import (
	"context"
	"fmt"
	"lhyim_server/common/list_query"
	"lhyim_server/common/models"
	"lhyim_server/lhyim_user/user_models"

	"lhyim_server/lhyim_user/user_api/internal/svc"
	"lhyim_server/lhyim_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	//先找所有的用户
	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.UserConfModel{
		Online: req.Online,
	}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Join:    "left join user_models um on um.id = user_conf_models.user_id",
		Preload: []string{"UserModel"},
		Where:   l.svcCtx.DB.Where("(user_conf_models.search_user <> 0 or user_conf_models.search_user is not null)  and (user_conf_models.search_user = 1 and um.id = ?) or (user_conf_models.search_user = 2 and (um.id = ? or um.nickname like  ?))  ", req.Key, req.Key, fmt.Sprintf("%%%s%%", req.Key)),
	})
	//查自己用户的好友列表
	var isfriend user_models.FriendModel
	isfriends := isfriend.Friends(l.svcCtx.DB, req.UserID)
	userMap := map[uint]bool{}
	for _, model := range isfriends {
		if model.SendUserID == req.UserID {
			userMap[model.RecvUserID] = true
		} else {
			userMap[model.SendUserID] = true
		}
	}
	list := make([]types.SearchInfo, 0)
	for _, friend := range friends {
		list = append(list, types.SearchInfo{
			UserID:   int64(friend.UserID),
			Nickname: friend.UserModel.Nickname,
			Avatar:   friend.UserModel.Avatar,
			Abstract: friend.UserModel.Abstract,
			IsFriend: userMap[friend.UserID],
		})
	}
	return &types.SearchResponse{
		List:  list,
		Count: count,
	}, nil
}
