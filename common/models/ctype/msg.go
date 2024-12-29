package ctype

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type MsgType int8

const (
	TextMsgType MsgType = iota + 1
	ImageMsgType
	VideoMsgType
	FileMsgType
	VoiceMsgType
	VoiceCallMsgType
	VideoCallMsgType
	WithdrawMsgType
	ReplyMsgType
	QuoteMsgType
	AtMsgType
	TipMsgType
)

type ImageMsg struct {
	Title string `json:"title"` //图片标题
	Src   string `json:"src"`   //图片
}
type VideoMsg struct {
	Title string `json:"title"` //视频标题
	Src   string `json:"src"`   //视频
	Time  int    `json:"time"`  //视频时长 单位秒
}
type FileMsg struct {
	Title string `json:"title"` //文件标题
	Src   string `json:"src"`   //文件
	Size  int64  `json:"size"`  //文件大小 单位字节

	Type string `json:"type"` //文件类型
}
type VoiceMsg struct {
	Src  string `json:"src"`  //视频
	Time int    `json:"time"` //语音时长 单位秒

}
type VoiceCallMsg struct {
	StartTime time.Time `json:"startTime"` //开始时间
	EndTime   time.Time `json:"endTime"`   //结束时间
	EndReason int8      `json:"endReason"` //结束原因 0发起方挂断 1接收方挂断 2网络原因 3未打通
}
type VideoCallMsg struct {
	StartTime time.Time `json:"startTime"` //开始时间
	EndTime   time.Time `json:"endTime"`   //结束时间
	EndReason int8      `json:"endReason"` //结束原因 0发起方挂断 1接收方挂断 2网络原因 3未打通
}
type WithdrawMsg struct {
	Content   string `json:"content"`             //撤回消息
	MsgID     uint   `json:"msgID"`               //撤回消息id
	OriginMsg *Msg   `json:"originMsg,omitempty"` //原消息

}
type ReplyMsg struct {
	MsgID         uint      `json:"msgID"`         //回复消息ID
	Content       string    `json:"content"`       //回复消息
	Msg           *Msg      `json:"msg,omitempty"` //原消息
	UserID        uint      `json:"userID"`        //被回复人id
	UserNickName  string    `json:"userNickName"`  //被回复人的昵称
	OriginMsgDate time.Time `json:"originMsgDate"`
}
type QuoteMsg struct {
	MsgID         uint      `json:"msgID"`        //回复消息ID
	Content       string    `json:"content"`      //回复消息
	Msg           *Msg      `json:"msg"`          //原消息
	UserID        uint      `json:"userID"`       //被回复人id
	UserNickName  string    `json:"userNickName"` //被回复人的昵称
	OriginMsgDate time.Time `json:"originMsgDate"`
}
type AtMsg struct {
	UserID  uint   `json:"userID"`  //用户ID
	Content string `json:"content"` //回复消息
	Msg     *Msg   `json:"msg"`     //原消息
}
type TextMsg struct {
	Content string `json:"content"`
	Src     string `json:"src"`
}
type TipMsg struct {
	Content string `json:"content"`
	Status  string `json:"status"` //error success info
}
type Msg struct {
	Type         MsgType       `json:"type"`                   //消息类型 1文字消息 2图片消息 3视频消息 4文件消息 5语音消息 6语言通话 7视频童话 8撤回消息 9回复消息 10引用消息
	TextMsg      *TextMsg      `json:"textMsg,omitempty"`      //文字消息
	ImageMsg     *ImageMsg     `json:"imageMsg,omitempty"`     //图片消息
	VideoMsg     *VideoMsg     `json:"videoMsg,omitempty"`     //视频消息
	FileMsg      *FileMsg      `json:"fileMsg,omitempty"`      //文件消息
	VoiceMsg     *VoiceMsg     `json:"voiceMsg,omitempty"`     //语音消息
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg,omitempty"` //语音通话
	VideoCallMsg *VideoCallMsg `json:"videoCallMsg,omitempty"` //视频通话
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg,omitempty"`  //撤回消息
	ReplyMsg     *ReplyMsg     `json:"replyMsg,omitempty"`     //回复消息
	QuoteMsg     *QuoteMsg     `json:"quoteMsg,omitempty"`     //引用消息
	AtMsg        *AtMsg        `json:"atMsg,omitempty"`        //at消息
	TipMsg       *TipMsg       `json:"tipMsg,omitempty"`       //提示消息，一般不入库
}

func (msg Msg) MsgPreview() string {
	switch msg.Type {
	case TextMsgType:
		return msg.TextMsg.Content
	case ImageMsgType:
		return "[图片]"
	case VideoMsgType:
		return "[视频]"
	case FileMsgType:
		return "[文件]"
	case VoiceMsgType:
		return "[语音]"
	case VoiceCallMsgType:
		return "[视频通话]"
	case VideoCallMsgType:
		return "[撤回消息]"
	case WithdrawMsgType:
		return "[回复消息]"
	case ReplyMsgType:
		return "[回复消息]"
	case QuoteMsgType:
		return "[引用消息]"
	default:
		panic("unhandled default case")

	}
	return "[未知消息]"
}
func (c *Msg) Scan(val interface{}) error {
	err := json.Unmarshal(val.([]byte), c)
	if err != nil {
		return err
	}
	if c.Type == WithdrawMsgType {
		if c.WithdrawMsg != nil {
			c.WithdrawMsg.OriginMsg = nil
		}
	}
	return nil
}
func (c Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
