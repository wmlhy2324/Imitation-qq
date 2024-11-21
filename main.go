package main

import (
	"flag"
	"fmt"
	"lhyim_server/core"
	"lhyim_server/lhyim_chat/chat_models"
	"lhyim_server/lhyim_file/file_model"
	"lhyim_server/lhyim_group/group_models"
	"lhyim_server/lhyim_user/user_models"
)

type Options struct {
	DB bool
}

func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()
	if opt.DB {
		db := core.InitGorm("root:112304@tcp(127.0.0.1:3306)/lhyim_server_db?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(&user_models.UserModel{},
			&user_models.UserConfModel{},
			&user_models.FriendModel{},
			&user_models.FriendVerifyModel{},
			&chat_models.ChatModel{},
			&group_models.GroupModel{},
			&group_models.GroupMemberModel{},
			&group_models.GroupVerifyModel{},
			&group_models.GroupMsgModel{},
			&chat_models.TopUserModel{},
			&chat_models.UserChatDeleteModel{}, //置顶用户表
			&file_model.FileModel{},
		)
		if err != nil {
			fmt.Println("表结构创建失败", err)
			return
		}
		fmt.Println("表结构创建成功")
	}
	fmt.Println("初始化成功")
}
