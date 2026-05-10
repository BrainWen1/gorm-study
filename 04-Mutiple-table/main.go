package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID         int
	Name       string
	UserDetail *UserDetail `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`  // 一对一关系; 级联删除
	VideoList  []*Video    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL"` // 一对多关系; 删除用户时，设置关联视频的 UserID 为 NULL
}

type UserDetail struct {
	ID      int
	Address string

	UserID int   // 外键
	User   *User // 一对一关系
}

type Video struct {
	ID    int
	Title string

	UserID int   // 外键
	User   *User // 一对多关系

	TagList []*Tag `gorm:"many2many:video_tags;constraint:OnDelete:CASCADE"` // 多对多关系
}

type Tag struct {
	ID   int
	Name string

	VideoList []*Video `gorm:"many2many:video_tags;"` // 多对多关系
}

// 初始化
var DB *gorm.DB

func init() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:123456@tcp(192.168.12.143:3306)/multi_tab"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}

func migrate() {
	// 自动迁移
	err := DB.AutoMigrate(
		&User{},
		&UserDetail{},
		&Video{},
		&Tag{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

func main() {
	migrate()
	// 一对一关系
	// 插入数据
	// DB.Create(&User{
	// 	Name: "Alice",
	// 	UserDetail: &UserDetail{
	// 		Address: "123 Main St",
	// 	},
	// })

	// DB.Create(&UserDetail{
	// 	UserID:  1,
	// 	Address: "456 Elm St",
	// })

	// 查询数据
	var user User
	DB.Preload("UserDetail").First(&user, 1) // Preload 预加载关联数据
	fmt.Println(user, user.UserDetail)

	var userDetail UserDetail
	DB.Preload("User").First(&userDetail, 1)
	fmt.Println(userDetail, userDetail.User)

	// 删除数据
	// 级联删除：删除用户时，自动删除关联的用户详情
	DB.Select("UserDetail").Delete(&User{}, 1)
}
