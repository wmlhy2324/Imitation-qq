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

func (chat ChatModel) MsgPreviewMethod() string {
	if chat.SystemMsg != nil {
		switch chat.SystemMsg.Type {
		case 1:
			return "[系统消息]- 该消息涉黄,已被系统拦截"
		case 2:
			return "[系统消息]- 该消息涉恐,已被系统拦截"
		case 3:
			return "[系统消息]- 该消息涉证,已被系统拦截"
		case 4:
			return "[系统消息]- 该消息涉及不当言论,已被系统拦截"
		}
		return "[系统消息]"
	}
	switch chat.Msg.Type {
	case 1:
		return chat.Msg.TextMsg.Content
	case 2:
		return "[图片]"
	case 3:
		return "[视频]"
	case 4:
		return "[文件]"
	case 5:
		return "[语音]"
	case 6:
		return "[视频通话]"
	case 7:
		return "[撤回消息]"
	case 8:
		return "[回复消息]"
	case 9:
		return "[引用消息]"
	case 10:
		return "[@消息]"

	}
	return "[未知消息]"
}
