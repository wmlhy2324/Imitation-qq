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
	"lhyim_server/lhyim_chat/chat_api/internal/svc"
	"lhyim_server/lhyim_chat/chat_api/internal/types"
	"lhyim_server/lhyim_chat/chat_models"
	"lhyim_server/lhyim_file/file_rpc/types/file_rpc"
	"lhyim_server/lhyim_user/user_models"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type UserWsInfo struct {
	UserInfo    user_models.UserModel      //用户的ws连接对象
	WsClientMap map[string]*websocket.Conn //这个用户管理的所有客户端
	Current     *websocket.Conn            //当前连接对象
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

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
		addr := conn.RemoteAddr().String()
		defer func() {
			conn.Close()
			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				//删除的是退出的那个ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				delete(UserOnlineWsMap, req.UserID)
				svcCtx.Redis.HDel("online", fmt.Sprintf("%d", req.UserID))
			}

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

		fmt.Println("addr = ", addr)
		userWsInfo, ok := UserOnlineWsMap[req.UserID]

		if !ok {
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
				Current: conn,
			}
			//说明用户是第一次来
			UserOnlineWsMap[req.UserID] = userWsInfo

		}
		_, ok1 := userWsInfo.WsClientMap[addr]
		if !ok1 {
			//表示用户二开以上
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
			//把当前连接对象更换
			UserOnlineWsMap[req.UserID].Current = conn
		}
		//把在线的用户存进redis
		svcCtx.Redis.HSet("online", fmt.Sprintf("%d", req.UserID), req.UserID)
		//查一下自己好友是不是上线
		//查一下自己的好友列表，返回用户id列表，看看在不在UserWsMap中，在的话就给自己发一个消息
		fmt.Println("是否开启了好友请求:", userInfo.UserConfModel.FriendOnline)

		//如果用户开启了好友上线提醒，查一下
		friendRes, err := svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		logx.Infof("好友上线:%s 用户id:%d", userInfo.Nickname, req.UserID)
		for _, info := range friendRes.FriendList {
			friend, ok := UserOnlineWsMap[uint(info.UserId)]
			text := fmt.Sprintf("好友%s上线", userInfo.Nickname)
			if ok {
				//好友上线了
				if friend.UserInfo.UserConfModel.FriendOnline {
					//friend.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("好友%s上线", userInfo.Nickname)))
					sendWsMapMsg(friend.WsClientMap, []byte(text))
				}
			}
		}

		for {
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				logx.Error(err)
				break
			}

			//发送消息
			var request ChatRequest

			err2 := json.Unmarshal(p, &request)
			if err2 != nil {
				//用户乱发消息
				SendTipMsg(conn, "参数解析失败")
				continue
			}
			if request.RevUserID != req.UserID {
				isFriendRes, err := svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
					User1: uint32(req.UserID),
					User2: uint32(request.RevUserID),
				})
				if err != nil {
					logx.Error(err)
					SendTipMsg(conn, "用户服务错误")
					continue
				}
				if !isFriendRes.IsFriend {
					SendTipMsg(conn, "你们还不是好友呢")
					continue
				}
			}
			if request.Msg.Type > 12 && request.Msg.Type < 1 {
				SendTipMsg(conn, "消息类型错误")
				continue
			}
			//判断是否是文件类型
			switch request.Msg.Type {
			case ctype.TextMsgType:
				if request.Msg.TextMsg == nil {
					SendTipMsg(conn, "请输入消息内容")
					continue
				}
				if request.Msg.TextMsg.Content == "" {
					SendTipMsg(conn, "请输入消息内容")
					continue
				}

			case ctype.FileMsgType:
				//文件类型就要请求文件服务
				nameList := strings.Split(request.Msg.FileMsg.Src, "/")
				if len(nameList) == 0 {
					SendTipMsg(conn, "请上传文件")
					continue
				}
				FileId := nameList[len(nameList)-1]
				fileResponse, err3 := svcCtx.FileRpc.FileInfo(context.Background(), &file_rpc.FileInfoRequest{
					FileId: FileId,
				})
				if err3 != nil {
					logx.Error(err3)
					SendTipMsg(conn, err3.Error())
					continue
				}

				request.Msg.FileMsg.Title = fileResponse.FileName
				request.Msg.FileMsg.Size = fileResponse.FileSize
				request.Msg.FileMsg.Type = fileResponse.FileType
			case ctype.WithdrawMsgType:
				//  撤回消息的id是必填的
				if request.Msg.WithdrawMsg.MsgID == 0 {
					SendTipMsg(conn, "撤回消息id未填写")
					continue
				}
				//自己撤回自己发的
				//找这个消息是谁发的
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.WithdrawMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "消息不存在")
					continue
				}
				if msgModel.SendUserID != req.UserID {
					SendTipMsg(conn, "只能撤回自己的消息")
					continue
				}
				//判断消息的时间，小于2分钟的才能撤回
				now := time.Now()
				subTime := now.Sub(msgModel.CreatedAt)
				if subTime >= time.Minute*2 {
					SendTipMsg(conn, "只能撤回两分钟以内的消息哦")
					continue
				}

				//撤回逻辑
				//收到撤回请求之后,服务端这边把原消息修改为撤回消息类型，并且记录原消息
				//然后通知前端的发收双方，重新拉取聊天记录
				var content = "撤回了一条消息"
				//前端判断
				if userInfo.UserConfModel.RecallMessage != nil {
					content = *userInfo.UserConfModel.RecallMessage
				}
				content = "你" + content
				originMsg := msgModel.Msg
				originMsg.WithdrawMsg = nil //这里可能会出现循环引用
				svcCtx.DB.Model(&msgModel).Updates(chat_models.ChatModel{
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
			case ctype.QuoteMsgType:
				if request.Msg.QuoteMsg.MsgID == 0 {
					SendTipMsg(conn, "引用消息id需要填写")
					continue
				}
				//找到这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "消息不存在")
					continue
				}
				//不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipMsg(conn, "该消息已撤回")
					continue
				}
				if !((msgModel.SendUserID == req.UserID && msgModel.RecvUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RecvUserID == req.UserID)) {
					SendTipMsg(conn, "只能回复自己或对方的消息")
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
			case ctype.ReplyMsgType:
				//先去校验
				if request.Msg.ReplyMsg.MsgID == 0 {
					SendTipMsg(conn, "回复消息id需要填写")
					continue
				}
				//找到这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipMsg(conn, "消息不存在")
					continue
				}
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipMsg(conn, "该消息已撤回")
					continue
				}
				if !((msgModel.SendUserID == req.UserID && msgModel.RecvUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RecvUserID == req.UserID)) {
					SendTipMsg(conn, "只能回复自己或对方的消息")
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
			}
			//先入库,有些不用入库
			MsgID := InsetMsgByChat(svcCtx.DB, request.RevUserID, req.UserID, request.Msg)

			//给发送双方都要发消息

			SendMagByUser(svcCtx, request.RevUserID, req.UserID, request.Msg, MsgID)

		}

	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"`
	Msg       ctype.Msg `json:"msg"`
}
type ChatResponse struct {
	ID       uint           `json:"id"`
	IsMe     bool           `json:"isMe"`
	SendUser ctype.UserInfo `json:"sendUserID"`
	RevUser  ctype.UserInfo `json:"revUserID"`
	Msg      ctype.Msg      `json:"msg"`
	Created  time.Time      `json:"created"`
}

func InsetMsgByChat(db *gorm.DB, RevUserID, SendUserID uint, msg ctype.Msg) (MsgId uint) {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		logx.Infof("撤回消息不入库")
		return
	}
	chatModel := chat_models.ChatModel{
		SendUserID: SendUserID,
		RecvUserID: RevUserID,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()

	err := db.Create(&chatModel).Error
	if err != nil {
		logx.Error(err)
		sendUser, ok := UserOnlineWsMap[SendUserID]
		if !ok {
			return
		}
		SendTipMsg(sendUser.Current, "消息保存失败")
	}
	return chatModel.ID
}

// 发给谁,谁发的，发什么
func SendMagByUser(svcCtx *svc.ServiceContext, RevUserID, SendUserID uint, msg ctype.Msg, MsgID uint) {
	resp := ChatResponse{
		ID:      MsgID,
		Msg:     msg,
		Created: time.Now(),
	}

	RevUser, ok1 := UserOnlineWsMap[RevUserID]
	SendUser, ok2 := UserOnlineWsMap[SendUserID]
	if ok1 && ok2 && SendUserID == RevUserID {
		resp.RevUser = ctype.UserInfo{
			ID:       RevUserID,
			Nickname: RevUser.UserInfo.Nickname,
			Avatar:   RevUser.UserInfo.Avatar,
		}
		resp.SendUser = ctype.UserInfo{
			ID:       SendUserID,
			Nickname: SendUser.UserInfo.Nickname,
			Avatar:   SendUser.UserInfo.Avatar,
		}
		byteDate, _ := json.Marshal(resp)
		//RevUser.Conn.WriteMessage(websocket.TextMessage, byteDate)
		sendWsMapMsg(RevUser.WsClientMap, byteDate)
		return

	}
	//如果接收者不在线，就要去拿接收者的用户信息
	if !ok1 {
		//接收者在线
		userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, RevUserID)
		if err5 != nil {
			logx.Error(err5)
			return
		}
		resp.RevUser = ctype.UserInfo{
			ID:       RevUserID,
			Nickname: userBaseInfo.Nickname,
			Avatar:   userBaseInfo.Avatar,
		}
	} else {
		resp.RevUser = ctype.UserInfo{
			ID:       RevUserID,
			Nickname: RevUser.UserInfo.Nickname,
			Avatar:   RevUser.UserInfo.Avatar,
		}
	}

	resp.SendUser = ctype.UserInfo{
		ID:       SendUserID,
		Nickname: SendUser.UserInfo.Nickname,
		Avatar:   SendUser.UserInfo.Avatar,
	}
	resp.IsMe = true
	byteDate, _ := json.Marshal(resp)
	//SendUser.Conn.WriteMessage(websocket.TextMessage, byteDate)
	sendWsMapMsg(SendUser.WsClientMap, byteDate)
	if ok1 {
		resp.IsMe = false
		byteDate, _ := json.Marshal(resp)
		//	RevUser.Conn.WriteMessage(websocket.TextMessage, byteDate)
		sendWsMapMsg(RevUser.WsClientMap, byteDate)
	}
}

// 发送错误提示的消息
func SendTipMsg(conn *websocket.Conn, msg string) {

	resp := ChatResponse{
		Msg: ctype.Msg{
			Type: ctype.TipMsgType,
			TipMsg: &ctype.TipMsg{
				Content: msg,
				Status:  "error",
			},
		},
		Created: time.Now(),
	}
	byteDate, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, byteDate)
}

// 给一组ws对象发消息
func sendWsMapMsg(wsmsg map[string]*websocket.Conn, byteData []byte) {
	for _, conn := range wsmsg {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
}
