## lhyim_server
### 1. 项目简介
####  服务拆分
####  网关
####  用户服务
####  群聊服务
####  文件服务
####  聊天服务
##### 在群组服务里面有个问题xx in （?)如果这里传一个切片，切片为空的话，整个查询都会受到影响，
##### 比如Where: l.svcCtx.DB.Where("id not in ?", msgIDList),在这个里面msgIDList为空，where出来会一直为false
