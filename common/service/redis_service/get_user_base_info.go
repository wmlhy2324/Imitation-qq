package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"lhyim_server/common/models/ctype"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"time"
)

func GetUserBaseInfo(client *redis.Client, UserRpc user_rpc.UsersClient, userID uint) (userInfo ctype.UserInfo, err error) {
	key := fmt.Sprintf("lim_server_user_%d", userID)
	str, err := client.Get(key).Result()
	if err != nil {
		//redis没找到
		res, err1 := UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(userID),
		})
		if err1 != nil {
			err = err1
			return
		}
		userInfo.ID = userID
		userInfo.Avatar = res.Avatar
		userInfo.Nickname = res.NickName
		byteUserInfo, _ := json.Marshal(userInfo)
		//设置进缓存
		client.Set(key, string(byteUserInfo), time.Hour)
		return userInfo, nil

	}
	err = json.Unmarshal([]byte(str), &userInfo)
	if err != nil {
		return
	}
	return
}
