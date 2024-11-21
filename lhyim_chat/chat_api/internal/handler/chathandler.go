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
	UserInfo user_models.UserModel
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
			svcCtx.Redis.HDel("online", fmt.Sprintf("%d", req.UserID))
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
			UserInfo: userInfo,

			Conn: conn,
		}
		UserWsMap[req.UserID] = userWsInfo
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
			friend, ok := UserWsMap[uint(info.UserId)]
			if ok {
				//好友上线了
				if friend.UserInfo.UserConfModel.FriendOnline {
					friend.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("好友%s上线", userInfo.Nickname)))
				}
			}
		}

		for {
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				logx.Error(err)
				break
			}
			fmt.Println(string(p), req.UserID)
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
			//判断是否是文件类型
			switch request.Msg.Type {
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
				fmt.Println(fileResponse)
			}
			//先入库
			InsetMsgByChat(svcCtx.DB, request.RevUserID, req.UserID, request.Msg)
			//看看目标用户在不在线
			SendMagByUser(request.RevUserID, req.UserID, request.Msg)

		}

	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"`
	Msg       ctype.Msg `json:"msg"`
}
type ChatResponse struct {
	SendUser ctype.UserInfo `json:"sendUserID"`
	RevUser  ctype.UserInfo `json:"revUserID"`
	Msg      ctype.Msg      `json:"msg"`
	Created  time.Time      `json:"created"`
}

func InsetMsgByChat(db *gorm.DB, RevUserID, SendUserID uint, msg ctype.Msg) {
	chatModel := chat_models.ChatModel{
		SendUserID: SendUserID,
		RecvUserID: RevUserID,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	err := db.Create(&chatModel)
	if err != nil {
		logx.Error(err)
		sendUser, ok := UserWsMap[SendUserID]
		if !ok {
			return
		}
		SendTipMsg(sendUser.Conn, "消息保存失败")
	}
}

// 发给谁,谁发的，发什么
func SendMagByUser(RevUserID, SendUserID uint, msg ctype.Msg) {

	RevUser, ok := UserWsMap[RevUserID]
	if !ok {
		return
	}
	sendUser, ok := UserWsMap[SendUserID]
	if !ok {
		return
	}

	resp := ChatResponse{
		RevUser: ctype.UserInfo{
			ID:       RevUserID,
			Nickname: RevUser.UserInfo.Nickname,
			Avatar:   RevUser.UserInfo.Avatar,
		},
		SendUser: ctype.UserInfo{
			ID:       SendUserID,
			Nickname: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		},
		Msg:     msg,
		Created: time.Now(),
	}
	byteDate, _ := json.Marshal(resp)
	RevUser.Conn.WriteMessage(websocket.TextMessage, byteDate)
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
