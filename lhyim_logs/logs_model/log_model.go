package logs_model

import "lhyim_server/common/models"

type LogModel struct {
	models.Model
	LogType      int8   `json:"logType"` //2操作日志 2运行日志
	Addr         string `json:"addr"`
	IP           string `json:"ip"`
	UserID       uint   `json:"userID"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	Level        string `json:"level"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Service      string `json:"service"` //记录微服务的名称
}
