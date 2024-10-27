package chat_models

import (
	"lhyim_server/common/models"
	"lhyim_server/common/models/ctype"
)

type ChatModel struct {
	models.Model
	SendUserID uint             `json:"sendUserID"`                 //发起验证方
	RecvUserID uint             `json:"recvUserID"`                 //接收验证方
	MsgType    uint             `json:"msgType"`                    //消息类型 1文字消息 2图片消息 3视频消息 4文件消息 5语音消息 6语言通话 7视频童话 8撤回消息 9回复消息 10引用消息
	MsgPreview string           `gorm:"size:256" json:"msgPreview"` //消息预览
	Msg        ctype.Msg        `json:"msg"`
	SystemMsg  *ctype.SystemMsg `json:"systemMsg"` //系统提示
}
