package ctype

import (
	"database/sql/driver"
	"encoding/json"
	"time"
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
	Content   string `json:"content"` //撤回消息
	OriginMsg *Msg   `json:"-"`       //原消息

}
type ReplyMsg struct {
	MsgID   uint   `json:"msgID"`   //回复消息ID
	Content string `json:"content"` //回复消息
	Msg     *Msg   `json:"msg"`     //原消息
}
type QuoteMsg struct {
	MsgID   uint   `json:"msgID"`   //回复消息ID
	Content string `json:"content"` //回复消息
	Msg     *Msg   `json:"msg"`     //原消息
}
type AtMsg struct {
	UserID  uint   `json:"userID"`  //用户ID
	Content string `json:"content"` //回复消息
	Msg     *Msg   `json:"msg"`     //原消息
}

type Msg struct {
	Type         uint          `json:"type"`         //消息类型 1文字消息 2图片消息 3视频消息 4文件消息 5语音消息 6语言通话 7视频童话 8撤回消息 9回复消息 10引用消息
	Content      *string       `json:"content"`      //文字消息
	ImageMsg     *ImageMsg     `json:"imageMsg"`     //图片消息
	VideoMsg     *VideoMsg     `json:"videoMsg"`     //视频消息
	FileMsg      *FileMsg      `json:"fileMsg"`      //文件消息
	VoiceMsg     *VoiceMsg     `json:"voiceMsg"`     //语音消息
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg"` //语音通话
	VideoCallMsg *VideoCallMsg `json:"videoCallMsg"` //视频通话
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg"`  //撤回消息
	ReplyMsg     *ReplyMsg     `json:"replyMsg"`     //回复消息
	QuoteMsg     *QuoteMsg     `json:"quoteMsg"`     //引用消息
	AtMsg        *AtMsg        `json:"atMsg"`        //at消息
}

func (c *Msg) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), c)
}
func (c Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
