package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(192.168.12.143:3306)/game"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println(db, "is connected")

	// 定义模型：表结构
	type Book struct {
		ID    uint `gorm:"primaryKey"`
		Name  string
		Price float64
	}

	// 自动迁移：创建表
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic(err)
	}

	// 创建记录：插入数据
	book := Book{Name: "CS:APP", Price: 60.0}
	err = db.Create(&book).Error
	if err != nil {
		panic(err)
	}

	// 查询记录
	var queriedBook Book
	err = db.First(&queriedBook, "name = ?", "CS:APP").Error // 查询单条记录
	if err != nil {
		panic(err)
	}
	fmt.Println(queriedBook)

	err = db.Take(&queriedBook, "id = ?", 2).Error // 查询任意一条记录
	if err == gorm.ErrRecordNotFound {
		fmt.Println("Record not found")
	}
	fmt.Println(queriedBook)
}
