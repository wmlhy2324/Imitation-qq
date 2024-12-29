package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"lhyim_server/common/models/ctype"
	"lhyim_server/common/response"
	"lhyim_server/lhyim_group/group_api/internal/svc"
	"lhyim_server/lhyim_group/group_api/internal/types"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserWsInfo struct {
	UserInfo    ctype.UserInfo
	WsClientMap map[string]*websocket.Conn //这个用户管理的客户端

}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

type ChatRequest struct {
	GroupID uint      `json:"groupID"`
	Msg     ctype.Msg `json:"msg"` //
}
type ChatResponse struct {
	UserID       uint          `json:"userID"`
	UserNickname string        `json:"userNickname"`
	UserAvatar   string        `json:"userAvatar"`
	Msg          ctype.Msg     `json:"msg"`
	ID           uint          `json:"id"`
	CreatedAt    time.Time     `json:"createdAt"`
	MsgType      ctype.MsgType `json:"msgType"`
	IsMe         bool          `json:"isMe"`
}

func groupChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupChatRequest
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
		addr := conn.RemoteAddr().String() //用户可能开多个客户端
		logx.Infof("用户建立ws连接 %s", addr)
		defer func() {
			conn.Close()
			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				//删除的是退出的那个ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				delete(UserOnlineWsMap, req.UserID)
				//svcCtx.Redis.HDel("online", fmt.Sprintf("%d", req.UserID))
			}

		}()
		baseInfoResponse, err := svcCtx.UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{UserId: uint32(req.UserID)})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		userInfo := ctype.UserInfo{
			ID:       req.UserID,
			Nickname: baseInfoResponse.NickName,
			Avatar:   baseInfoResponse.Avatar,
		}
		userWsInfo, ok := UserOnlineWsMap[req.UserID]

		if !ok {
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
			}
			//说明用户是第一次来
			UserOnlineWsMap[req.UserID] = userWsInfo

		}
		_, ok1 := userWsInfo.WsClientMap[addr]
		if !ok1 {
			//表示用户二开以上
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
		}
		for {
			//消息类型，消息
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				break
			}
			var request ChatRequest
			err = json.Unmarshal(p, &request)
			if err != nil {
				SendTipMsg(conn, "参数解析失败")
				continue
			}
			//判断自己是不是群成员
			var member group_models.GroupMemberModel
			err = svcCtx.DB.Take(&member, "user_id = ? and group_id = ?", req.UserID, request.GroupID).Error
			if err != nil {
				SendTipMsg(conn, "你还不是这个群的成员呢")
				//自己不是群成员
				continue
			}
			msgID := InsetMsg(svcCtx.DB, conn, request.GroupID, req.UserID, request.Msg)
			//遍历这个用户列表，去找ws的客户端
			sendGroupOnlineUserMsg(svcCtx.DB, request.GroupID, req.UserID, request.Msg, msgID)
			fmt.Println(string(p))
		}
	}
}
func getOnlineUserIDList() (userOnlineIDList []uint) {
	for u := range UserOnlineWsMap {
		userOnlineIDList = append(userOnlineIDList, u)
	}
	return userOnlineIDList
}

func SendTipMsg(conn *websocket.Conn, msg string) {

	resp := ChatResponse{
		Msg: ctype.Msg{
			Type: ctype.TipMsgType,
			TipMsg: &ctype.TipMsg{
				Content: msg,
				Status:  "error",
			},
		},
		CreatedAt: time.Now(),
	}
	byteDate, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, byteDate)
}

// 这个群的用户发消息
func sendGroupOnlineUserMsg(db *gorm.DB, groupID uint, userID uint, msg ctype.Msg, msgID uint) {
	userOnlineIDList := getOnlineUserIDList()
	//去查这个群的成员，并且在线
	var groupMemberOnlineIDList []uint
	db.Model(&group_models.GroupMemberModel{}).Where("group_id = ? and user_id in (?)", groupID, userOnlineIDList).
		Select("user_id").Scan(&groupMemberOnlineIDList)
	//构造响应
	var chatResponse = ChatResponse{
		UserID:    userID,
		Msg:       msg,
		ID:        msgID,
		MsgType:   msg.Type,
		CreatedAt: time.Now(),
	}
	wsInfo, ok := UserOnlineWsMap[userID]
	if ok {
		chatResponse.UserNickname = wsInfo.UserInfo.Nickname
		chatResponse.UserAvatar = wsInfo.UserInfo.Avatar
	}

	for _, id := range groupMemberOnlineIDList {
		wsUserInfo, ok2 := UserOnlineWsMap[id]
		if !ok2 {
			continue
		}
		chatResponse.IsMe = false
		if wsUserInfo.UserInfo.ID == userID {
			chatResponse.IsMe = true
		}
		byteDate, _ := json.Marshal(chatResponse)
		//多端发送
		for _, w2 := range wsUserInfo.WsClientMap {
			w2.WriteMessage(websocket.TextMessage, byteDate)
		}
	}
}
func InsetMsg(db *gorm.DB, conn *websocket.Conn, groupId uint, userID uint, msg ctype.Msg) uint {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		logx.Infof("撤回消息不入库")
		return 0
	}
	groupMsg := group_models.GroupMsgModel{
		GroupID:    groupId,
		SendUserID: userID,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	groupMsg.MsgPreview = groupMsg.MsgPreviewMethod()

	err := db.Create(&groupMsg).Error
	if err != nil {
		logx.Error(err)

		SendTipMsg(conn, "消息保存失败")
		return 0
	}
	return groupMsg.ID
}
