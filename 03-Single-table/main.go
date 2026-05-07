package main

import (
	"fmt"
	"gorm-study/03-Single-table/global"

	"gorm.io/gorm"
)

type UserModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50); not null"`
	Age  int    `gorm:"default:18"`
}

// 指定表名
// func (UserModel) TableName() string {
// 	return "users"
// }

func migrate() {
	err := global.DB.AutoMigrate(&UserModel{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func selectAll() {
	var users []UserModel
	err := global.DB.Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

// 钩子函数
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Hook function is called")
	return nil
}

func insert(users []UserModel) {
	err := global.DB.Create(&users).Error
	if err != nil {
		panic("failed to create users")
	}
}

func main() {
	// 连接数据库
	global.Connect()

	// 自动迁移
	migrate()

	// 查询记录
	var users []UserModel
	global.DB.Find(&users)
	fmt.Println(users)

	// // 插入数据
	// user1 := UserModel{Name: "Alice", Age: 32} // 回填插入
	// insert([]UserModel{user1})

	// selectAll()

	// // 插入多条数据
	// users = []UserModel{
	// 	{Name: "Bob", Age: 28},
	// 	{Name: "Charlie", Age: 25},
	// }
	// insert(users)

	// selectAll()

	// 条件查询
	err := global.DB.Find(&users, "age > ?", 30).Error // ? 作为占位符，防止 SQL 注入
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)

	// 单条记录查询: First、Last、Take, 实际开发中基本是使用 First
	var user UserModel
	err = global.DB.Where("name = ?", "Alice").First(&user).Error
	if err != nil {
		panic("failed to query user")
	}
	fmt.Println(user)

	// 打印原始 SQL 语句: Debug()
	err = global.DB.Debug().Where("name = ?", "Alice").First(&user).Error
	if err != nil {
		panic("failed to query user")
	}
	fmt.Println(user)
}
