package mqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"lhyim_server/lhyim_logs/logs_api/internal/svc"
	"lhyim_server/lhyim_logs/logs_model"
	"lhyim_server/lhyim_user/user_rpc/types/user_rpc"
	"sync"
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
	//判断是不是运行日志
	if info.LogType == 3 {
		//先查一下今天这个服务有没有日志，有的话就更新，没有就创建
		var logModel = logs_model.LogModel{}
		mutex := sync.Mutex{}
		mutex.Lock()
		err = l.svcCtx.DB.Take(&logModel, "log_type = ? and service = ? and to_days(created_at) = to_days(now())", 3, info.Service).Error
		mutex.Unlock()
		if err == nil {
			//找到了
			l.svcCtx.DB.Model(&logModel).Update("content", logModel.Content+"\n"+info.Content)
			logx.Infof("运行日志 %s 更新成功", req.Title)
			return nil
		}
	}
	logx.Infof("LogEvent key :%s , val :%s", key, val)
	mutex := sync.Mutex{}
	mutex.Lock()
	err = l.svcCtx.DB.Create(&info).Error
	mutex.Unlock()
	if err != nil {
		logx.Error(err)
		return nil
	}
	return nil
}
