package logic

import (
	"context"
	"errors"
	"fmt"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"lhyim_server/utils/set"
	"strings"
	"time"

	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupCreateLogic) GroupCreate(req *types.GroupCreateRequest) (resp *types.GroupCreateResponse, err error) {
	// todo: add your logic here and delete this line
	var groupModel = group_models.GroupModel{
		Creator:      req.UserID,
		IsSearch:     false,
		Abstract:     fmt.Sprintf("本群创建于%s:群主很懒,什么都没有留下", time.Now().Format("2006-01-02 15:04:05")),
		Verification: 2,
		Size:         50,
	}
	var groupUserList = []uint{req.UserID}
	switch req.Mode {
	case 1: //直接创建
		if req.Name == "" {
			return nil, errors.New("群名未填写")
		}
		if req.Size >= 1000 {
			return nil, errors.New("群规模错误")
		}
		groupModel.Title = req.Name
		groupModel.Size = req.Size
		groupModel.IsSearch = req.IsSearch
	case 2:
		if len(req.UserIDList) == 0 {
			return nil, errors.New("没有要创建的好友")
		}
		//群名最大32
		var userIDList = []uint32{
			uint32(req.UserID),
		}
		for _, u := range req.UserIDList {
			userIDList = append(userIDList, uint32(u))
			groupUserList = append(groupUserList, u)
		}
		//选择的用户id必须都是我的好友
		friendListResponse, err := l.svcCtx.UserRpc.FriendList(l.ctx, &user_rpc.FriendListRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("好友信息或许错误")
		}
		var friendIDList []uint
		for _, i2 := range friendListResponse.FriendList {
			friendIDList = append(friendIDList, uint(i2.UserId))
		}
		//判断他们两个是不是不一致的
		fmt.Println(req.UserIDList)
		fmt.Println(friendIDList)
		slice := set.Difference(req.UserIDList, friendIDList)
		if len(slice) != 0 {
			return nil, errors.New("选择的好友中不是你的好友")
		}
		userListResponse, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoRequest{
			UserIdList: userIDList,
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("创建获取用户信息错误")
		}

		//去算这个昵称的长度，多久会大于32
		var nameList []string
		for _, info := range userListResponse.UserInfo {
			if len(strings.Join(nameList, ",")) >= 29 {
				break
			}
			nameList = append(nameList, info.NickName)
		}
		groupModel.Title = strings.Join(nameList, ",") + "的群聊"
	default:
		return nil, errors.New("不支持的模式")
	}
	//群头像 1 默认头像 2 文字头像
	groupModel.Avatar = string([]rune(groupModel.Title)[0])
	err1 := l.svcCtx.DB.Create(&groupModel).Error
	if err1 != nil {
		logx.Error(err)
		return nil, errors.New("创建群聊失败")
	}
	//
	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {
		memberModel := group_models.GroupMemberModel{
			GroupID: groupModel.ID,
			UserID:  u,
			Role:    3,
		}
		if i == 0 {
			memberModel.Role = 1
		}
		members = append(members, memberModel)
	}
	err = l.svcCtx.DB.Create(&members).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群成员添加失败")
	}
	return
}
