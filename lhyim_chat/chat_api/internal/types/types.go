// Code generated by goctl. DO NOT EDIT.
package types

type ChatHistoryRequest struct {
	UserID   uint `header:"User-ID"`
	Page     int  `form:"page,optional"`
	Limit    int  `form:"limit,optional"`
	FriendID uint `form:"friendId"`
}

type ChatHistoryResponse struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"userId"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	CreateAt string `json:"createAt"`
}

type ChatSession struct {
	UserID     uint   `json:"userId"`
	Avatar     string `json:"avatar"`
	Nickname   string `json:"nickname"`
	CreateAt   string `json:"createAt"`   //消息时间
	MsgPreview string `json:"msgPreview"` //消息预览
	IsTop      bool   `json:"isTop"`      //是否置顶
}

type ChatSessionRequest struct {
	UserID uint `header:"User-ID"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
	Key    int  `form:"key,optional"`
}

type ChatSessionResponse struct {
	List  []ChatSession `json:"list"`
	Count int64         `json:"count"`
}

type UserTopRequest struct {
	UserID   uint `header:"User-ID"`
	FriendID uint `json:"friendId"`
}

type UserTopResponse struct {
}
