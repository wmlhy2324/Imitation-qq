package chat_models

// 用户删除聊天记录的表
type UserChatDeleteModel struct {
	UserID uint `json:"userID"`
	ChatID uint `json:"chatID"`
}
