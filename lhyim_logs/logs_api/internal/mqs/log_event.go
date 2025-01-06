package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/lhyim_logs/logs_api/internal/svc"
	"lhyim_server/lhyim_logs/logs_model"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
)

type LogEvent struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *LogEvent {
	return &LogEvent{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type Request struct {
	LogType int8   `json:"logType"` //2操作日志 2运行日志
	IP      string `json:"ip"`
	UserID  uint   `json:"userID"`
	Level   string `json:"level"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Service string `json:"service"` //记录微服务的名称
}

func (l *LogEvent) Consume(ctx context.Context, key, val string) error {
	var req Request
	err := json.Unmarshal([]byte(val), &req)
	if err != nil {
		logx.Error("json解析错误val = %s err = %s", val, err.Error())
		return nil
	}
	//查ip对应的地址
	var info = logs_model.LogModel{
		LogType: req.LogType,
		IP:      req.IP,
		UserID:  req.UserID,
		Addr:    "内网地址",
		Level:   req.Level,
		Title:   req.Title,
		Content: req.Content,
		Service: req.Service,
	}
	//调用户基础方法获取用户昵称
	baseInfo, err := l.svcCtx.UserRpc.UserBaseInfo(l.ctx, &user_rpc.UserBaseInfoRequest{UserId: uint32(req.UserID)})
	if err == nil {
		info.UserNickname = baseInfo.NickName
		info.UserAvatar = baseInfo.Avatar
	}
	logx.Infof("LogEvent key :%s , val :%s", key, val)
	err = l.svcCtx.DB.Create(&info).Error
	if err != nil {
		logx.Error(err)
		return nil
	}
	return nil
}
