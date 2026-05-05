package main

import (
	"gorm-study/03-Single-table/global"
)

type UserModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50); not null; unique"`
	Age  int    `gorm:"default:18"`
}

func migrate() {
	err := global.DB.AutoMigrate(&UserModel{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func main() {
	// 连接数据库
	global.Connect()

	// 自动迁移
	migrate()
}
