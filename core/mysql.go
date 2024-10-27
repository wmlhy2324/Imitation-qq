package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func InitGorm(MysqlDateSource string) *gorm.DB {

	var mysqlLogger logger.Interface
	db, err := gorm.Open(mysql.Open(MysqlDateSource), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("mysql connect error: %v", err))
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)               //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)              //最多可容纳数
	sqlDB.SetConnMaxLifetime(time.Hour * 4) //复用时间
	return db
}
