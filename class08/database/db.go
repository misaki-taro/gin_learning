package database

import (
	"class08/config"
	"class08/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 连接数据库
func InitDB() (*gorm.DB, error) {
	// dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	config.Init()
	db, err := gorm.Open(mysql.Open(config.Cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 建表
	err = db.AutoMigrate(model.User{}, model.Todo{})
	if err != nil {
		panic("failed to connect database")
	}

	return db, nil
}
