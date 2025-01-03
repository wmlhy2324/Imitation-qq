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
	"lhyim_server/common/service/redis_service"
	"lhyim_server/lhyim_chat/chat_models"
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
	UserID         uint          `json:"userID"`
	UserNickname   string        `json:"userNickname"`
	UserAvatar     string        `json:"userAvatar"`
	Msg            ctype.Msg     `json:"msg"`
	ID             uint          `json:"id"`
	CreatedAt      time.Time     `json:"createdAt"`
	MsgType        ctype.MsgType `json:"msgType"`
	IsMe           bool          `json:"isMe"`
	MemberNickname string        `json:"memberNickname"`
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
			err = svcCtx.DB.Preload("GroupModel").Take(&member, "user_id = ? and group_id = ?", req.UserID, request.GroupID).Error
			if err != nil {
				SendTipMsg(conn, "你还不是这个群的成员呢")
				//自己不是群成员
				continue
			}
			if member.GroupModel.IsProhibition {
				//开启了全员禁言
				SendTipMsg(conn, "当前群正在全员禁言中")
				continue
			}
			//我是不是被禁言了
			if member.GetProhibitionTime(svcCtx.Redis, svcCtx.DB) != nil {
				SendTipMsg(conn, "当前用户禁言中")
				continue
			}
			switch request.Msg.Type {
			case ctype.WithdrawMsgType: //撤回消息
				//拿到消息
				withdrawMsg := request.Msg.WithdrawMsg
				if withdrawMsg == nil {
					SendTipMsg(conn, "撤回消息的格式错误")
					continue
				}
				if withdrawMsg.MsgID == 0 {
					SendTipMsg(conn, "撤回消息id为空")
					continue
					//找到消息
				}
				var groupMsg group_models.GroupMsgModel
				err = svcCtx.DB.Take(&groupMsg, "group_id = ? and id = ?", request.GroupID, withdrawMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "原消息不存在")
					continue
				}
				//原消息不能是撤回消息
				if groupMsg.MsgType == ctype.WithdrawMsgType {
					SendTipMsg(conn, "该消息已撤回")
					continue
				}
				//拿我在这个群的角色
				if member.Role == 3 {
					if req.UserID != groupMsg.SendUserID {
						SendTipMsg(conn, "普通用户只能撤回自己的消息")
						continue
					}
					if req.UserID == groupMsg.SendUserID {
						now := time.Now()
						if now.Sub(groupMsg.CreatedAt) > 2*time.Minute {
							SendTipMsg(conn, "只能撤回两分钟以内的消息")
							continue
						}
					}
				}
				//如果自己撤回自己的

				//如果是群主或者管理员，可以撤回自己的，没有时间限制
				//查这个消息的用户在这个群里的角色
				var msgUserRole int8 = 3
				svcCtx.DB.Model(&group_models.GroupMemberModel{}).
					Where("group_id = ? and user_id = ?", request.GroupID, groupMsg.SendUserID).Select("role").Scan(&msgUserRole)

				//用户退群的情况

				if member.Role == 2 {
					if msgUserRole == 1 || (msgUserRole == 2 && req.UserID != groupMsg.SendUserID) {
						SendTipMsg(conn, "管理员只能撤回自己或者普通用户的消息")
						continue
					}
				}
				//消息可以撤回了
				//修改原消息
				var content = "撤回了一条消息"

				content = "你" + content
				originMsg := groupMsg.Msg
				originMsg.WithdrawMsg = nil //这里可能会出现循环引用
				//self
				request.Msg.WithdrawMsg.Content = content
				svcCtx.DB.Model(&groupMsg).Updates(chat_models.ChatModel{
					MsgPreview: "[撤回消息]-" + content,
					MsgType:    ctype.WithdrawMsgType,
					Msg: ctype.Msg{
						Type: ctype.WithdrawMsgType,
						WithdrawMsg: &ctype.WithdrawMsg{
							Content:   content,
							MsgID:     request.Msg.WithdrawMsg.MsgID,
							OriginMsg: &originMsg,
						},
					},
				})
			case ctype.ReplyMsgType:
				//先去校验
				if request.Msg.ReplyMsg.MsgID == 0 || request.Msg.ReplyMsg == nil {
					SendTipMsg(conn, "回复消息id需要填写")
					continue
				}
				//找到这个原消息
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id = ? and id = ?", request.GroupID, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "消息不存在")
					continue
				}
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipMsg(conn, "该消息已撤回")
					continue
				}

				userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, req.UserID)
				if err5 != nil {
					logx.Error(err5)
					SendTipMsg(conn, err5.Error())
					return
				}
				request.Msg.ReplyMsg.Msg = &msgModel.Msg
				request.Msg.ReplyMsg.UserID = msgModel.SendUserID
				request.Msg.ReplyMsg.UserNickName = userBaseInfo.Nickname
				request.Msg.ReplyMsg.OriginMsgDate = msgModel.CreatedAt
			case ctype.QuoteMsgType:
				//先去校验
				if request.Msg.QuoteMsg.MsgID == 0 || request.Msg.QuoteMsg == nil {
					SendTipMsg(conn, "回复消息id需要填写")
					continue
				}
				//找到这个原消息
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id = ? and id = ?", request.GroupID, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "消息不存在")
					continue
				}
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipMsg(conn, "该消息已撤回")
					continue
				}

				userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, req.UserID)
				if err5 != nil {
					logx.Error(err5)
					SendTipMsg(conn, err5.Error())
					return
				}
				request.Msg.QuoteMsg.Msg = &msgModel.Msg
				request.Msg.QuoteMsg.UserID = msgModel.SendUserID
				request.Msg.QuoteMsg.UserNickName = userBaseInfo.Nickname
				request.Msg.QuoteMsg.OriginMsgDate = msgModel.CreatedAt
			}

			msgID := InsetMsg(svcCtx.DB, conn, member, request.Msg)
			//遍历这个用户列表，去找ws的客户端
			sendGroupOnlineUserMsg(svcCtx.DB, member, request.Msg, msgID)
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
func sendGroupOnlineUserMsg(db *gorm.DB, member group_models.GroupMemberModel, msg ctype.Msg, msgID uint) {
	userOnlineIDList := getOnlineUserIDList()
	//去查这个群的成员，并且在线
	var groupMemberOnlineIDList []uint
	db.Model(&group_models.GroupMemberModel{}).Where("group_id = ? and user_id in (?)", member.GroupID, userOnlineIDList).
		Select("user_id").Scan(&groupMemberOnlineIDList)
	//构造响应
	var chatResponse = ChatResponse{
		UserID:         member.UserID,
		Msg:            msg,
		ID:             msgID,
		MsgType:        msg.Type,
		CreatedAt:      time.Now(),
		MemberNickname: member.MemberNickname,
	}

	wsInfo, ok := UserOnlineWsMap[member.UserID]
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
		if wsUserInfo.UserInfo.ID == member.UserID {
			chatResponse.IsMe = true
		}
		byteDate, _ := json.Marshal(chatResponse)
		//多端发送
		for _, w2 := range wsUserInfo.WsClientMap {
			w2.WriteMessage(websocket.TextMessage, byteDate)
		}
	}
}
func InsetMsg(db *gorm.DB, conn *websocket.Conn, member group_models.GroupMemberModel, msg ctype.Msg) uint {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		logx.Infof("撤回消息不入库")
		return 0
	}
	groupMsg := group_models.GroupMsgModel{
		GroupID:    member.GroupID,
		SendUserID: member.UserID,
		MemberID:   member.ID,
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
