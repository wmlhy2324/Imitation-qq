package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/common/response"
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"lhyim_server/lhyim_user/user_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserInfo struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	UserID   uint   `json:"userID"`
}
type UserWsInfo struct {
	UserInfo UserInfo
	Conn     *websocket.Conn //用户的ws连接对象
}

var UserWsMap = map[uint]UserWsInfo{}

func chatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		var upGrader = websocket.Upgrader{

			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		defer func() {
			conn.Close()
			delete(UserWsMap, req.UserID)
		}()
		res, err := svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
			UserId: uint32(req.UserID)},
		)
		if err != nil {
			logx.Error(err)
			return
		}
		var userInfo user_models.UserModel
		err = json.Unmarshal(res.Data, &userInfo)
		if err != nil {
			logx.Error(err)
			return
		}
		var userWsInfo = UserWsInfo{
			UserInfo: UserInfo{
				UserID:   req.UserID,
				Avatar:   userInfo.Avatar,
				Nickname: userInfo.Nickname,
			},
			Conn: conn,
		}
		UserWsMap[req.UserID] = userWsInfo
		//查一下自己好友是不是上线
		//查一下自己的好友列表，返回用户id列表，看看在不在UserWsMap中，在的话就给自己发一个消息
		if userInfo.UserConfModel.FriendOnline {
			//如果用户开启了好友上线提醒，查一下
			friendRes, err := svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
				UserId: uint32(req.UserID),
			})
			if err != nil {
				logx.Error(err)
				response.Response(r, w, nil, err)
				return
			}
			for _, info := range friendRes.FriendList {
				friend, ok := UserWsMap[uint(info.UserId)]
				if ok {
					//好友上线了
					conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("好友%s上线", friend.UserInfo.Nickname)))
				}
			}
		}
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				logx.Error(err)
				break
			}
			fmt.Println(string(p), req.UserID)
			//发送消息
			conn.WriteMessage(websocket.TextMessage, []byte("xxx"))
		}

	}
}
