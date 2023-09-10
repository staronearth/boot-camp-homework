package ioc

import (
	"boot-camp-homework/webook/config"
	"boot-camp-homework/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN), &gorm.Config{})
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic相当于整个goroutine结束
		// 一旦初始化过程出错,应用就不要启动了
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
