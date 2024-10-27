// Code generated by goctl. DO NOT EDIT.
package types

type AddFriendRequest struct {
	UserID               uint                  `header:"User-ID"`
	FriendID             uint                  `json:"friendId"`
	Verify               string                `json:"verify,optional"` //验证消息
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion"`
}

type AddFriendResponse struct {
}

type FriendDeleteRequest struct {
	UserID   uint `header:"User-ID"`
	FriendID uint `json:"friendId"`
}

type FriendDeleteResponse struct {
}

type FriendInfoResponse struct {
	UserID   int64  `json:"UserID"`
	Nickname string `gorm:"size:32" json:"nickname"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	Abstract string `gorm:"size:128" json:"abstract"`
	Notice   string `json:"notice"`
}

type FriendListRequest struct {
	UserID uint `header:"User-ID"`
	Role   int8 `header:"Role"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type FriendListResponse struct {
	List  []FriendInfoResponse `json:"list"`
	Count int                  `json:"count"`
}

type FriendNoticeUpdateRequest struct {
	UserID   uint   `header:"User-ID"`
	FriendID uint   `json:"friendID"`
	Notice   string `json:"notice"`
}

type FriendNoticeUpdateResponse struct {
}

type FriendValidInfo struct {
	UserID              uint                  `json:"userID"`
	Nickname            string                `json:"nickname"`
	Avatar              string                `json:"avatar"`
	AddtionalMessage    string                `json:"addtionalMessages"`
	VerficationQuestion *VerificationQuestion `json:"verficationQuestion"`
	Status              int8                  `json:"status"`       //状态0 未操作 1 同意 2 拒绝 3 忽略
	Verification        int8                  `json:"verification"` //好友验证
	ID                  uint                  `json:"id"`           //验证记录的id
	Flag                string                `json:"flag"`         //send我是发起方 recv 我是接收方
}

type FriendValidRequest struct {
	UserID uint `header:"User-ID"`
	Page   int  `form:"page,optional"`
	Limit  int  `form:"limit,optional"`
}

type FriendValidResponse struct {
	List  []FriendValidInfo `json:"list"`
	Count int64             `json:"count"`
}

type FriendValidStatusRequest struct {
	UserID   uint `header:"User-ID"`
	VerifyID uint `json:"verifyId"`
	Status   int8 `json:"status"` //状态
}

type FriendValidStatusResponse struct {
}

type SearchInfo struct {
	UserID   int64  `json:"UserID"`
	Nickname string `gorm:"size:32" json:"nickname"`
	Avatar   string `gorm:"size:256" json:"avatar"`
	Abstract string `gorm:"size:128" json:"abstract"`
	IsFriend bool   `json:"isFriend"`
}

type SearchRequest struct {
	UserID uint   `header:"User-ID"`
	Key    string `form:"key"`    //用户id和昵称
	Online bool   `form:"online"` //在线
	Page   int    `form:"page,optional"`
	Limit  int    `form:"limit,optional"`
}

type SearchResponse struct {
	List  []SearchInfo `json:"list"`
	Count int64        `json:"count"`
}

type UserInfoRequest struct {
	UserID uint `header:"User-ID"`
	Role   int8 `header:"Role"`
}

type UserInfoResponse struct {
	UserID               int64                `json:"UserID"`
	Nickname             string               `gorm:"size:32" json:"nickname"`
	Avatar               string               `gorm:"size:256" json:"avatar"`
	Abstract             string               `gorm:"size:128" json:"abstract"`
	ReCallMessage        *string              `gorm:"size:128" json:"reCallMessage"`
	FriendOnline         bool                 `json:"friendOnline"`
	Sound                bool                 `json:"sound"`
	SecureLink           bool                 `json:"secureLink"`
	SavePwd              bool                 `json:"savePwd"`
	SearchUser           int8                 `json:"searchUser"`
	Verification         int8                 `json:"verification"`
	VerificationQuestion VerificationQuestion `json:"verificationQuestion"`
}

type UserInfoUpdateRequest struct {
	UserID               uint                  `header:"User-ID"`
	Nickname             *string               `json:"nickname,optional" user:"nickname"`
	Abstract             *string               `json:"abstract,optional" user:"abstract"`
	Avatar               *string               `json:"avatar,optional" user:"avatar"`
	ReCallMessage        *string               `json:"recallMessage,optional" user_conf:"recall_message"`
	FriendOnline         *bool                 `json:"friendOnline,optional" user_conf:"friend_online"`
	Sound                *bool                 `json:"sound,optional" user_conf:"sound"`
	SecureLink           *bool                 `json:"secureLink,optional" user_conf:"secure_link"`
	SavePwd              *bool                 `json:"savePwd,optional" user_conf:"save_pwd"`
	SearchUser           *int8                 `json:"searchUser,optional" user_conf:"search_user"`
	Verification         *int8                 `json:"verification,optional" user_conf:"verification"`
	VerificationQuestion *VerificationQuestion `json:"verificationQuestion,optional" user_conf:"verification_question"`
}

type UserInfoUpdateResponse struct {
}

type UserValidRequest struct {
	UserID   uint `header:"User-ID"`
	FriendID uint `json:"friendId"`
}

type UserValidResponse struct {
	Verification         int8                 `json:"verification"`
	VerificationQuestion VerificationQuestion `json:"verificationQuestion"` //答案不反悔
}

type VerificationQuestion struct {
	Problem1 *string `json:"problem1,optional" user_conf:"problem1"`
	Problem2 *string `json:"problem2,optional" user_conf:"problem2"`
	Problem3 *string `json:"problem3,optional" user_conf:"problem3"`
	Answer1  *string `json:"answer1,optional" user_conf:"answer1"`
	Answer2  *string `json:"answer2,optional" user_conf:"answer2"`
	Answer3  *string `json:"answer3,optional" user_conf:"answer3"`
}

type FriendInfoRequest struct {
	UserID   uint `header:"User-ID"`
	Role     int8 `header:"Role"`
	FriendID uint `form:"friendID"` //form 表示是form表单提交的数据
}
