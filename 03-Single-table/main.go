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
	fmt.Println("Create's Hook function is called")
	return nil
}

func (u *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("Update's Hook function is called")
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

	// 更新数据
	fmt.Println("------------there is a split line------------")
	// Save : 根据主键更新，如果主键不存在则插入新记录
	selectAll()
	user.Age = 33
	err = global.DB.Save(&user).Error
	if err != nil {
		panic("failed to update user")
	}
	fmt.Println(user)
	selectAll()

	// Update : 更新单个字段
	err = global.DB.Model(&user).Update("age", 34).Error // Model() 指定要更新的模型，Update() 更新单个字段
	if err != nil {
		panic("failed to update user")
	}
	fmt.Println(user)
	selectAll()

	// UpdateColumn : 与update相比，只是不会触发钩子函数和更新字段的时间戳
	err = global.DB.Model(&user).UpdateColumn("age", 35).Error
	if err != nil {
		panic("failed to update user")
	}
	fmt.Println(user)
	selectAll()

	// Updates : 更新多个字段
	err = global.DB.Model(&user).Updates(UserModel{Age: 36, Name: "Alice Updated"}).Error
	if err != nil {
		panic("failed to update user")
	}
	fmt.Println(user)
	selectAll()

	// Expr : 使用表达式更新字段值
	err = global.DB.Model(&user).Update("age", gorm.Expr("age + ?", 1)).Error
	if err != nil {
		panic("failed to update user")
	}
	fmt.Println(user)
	selectAll()

	// 删除数据
	fmt.Println("------------there is a split line------------")
	user.ID = 3 // 删除 ID 为 3 的记录
	err = global.DB.Delete(&user).Error
	if err != nil {
		panic("failed to delete user")
	}
	selectAll()

	insert([]UserModel{{Name: "Alice", Age: 33}})
	selectAll()
	err = global.DB.Where("name = ? and age = ?", "Alice", 33).Delete(&UserModel{}).Error
	if err != nil {
		panic("failed to delete user")
	}
	selectAll()

	err = global.DB.Delete(&UserModel{}, []int{4, 5}).Error // 批量删除数据
	if err != nil {
		panic("failed to delete users")
	}
	selectAll()

	// 软删除 : 需要在模型中添加 DeletedAt 字段
	// type UserModel struct {
	// 	ID        uint           `gorm:"primaryKey"`
	// 	Name      string         `gorm:"type:varchar(50); not null"`
	// 	Age       int            `gorm:"default:18"`
	// 	DeletedAt gorm.DeletedAt `gorm:"index"` // 添加 DeletedAt 字段，启用软删除功能
	// }
	// 删除记录时，实际上是将 DeletedAt 字段设置为当前时间，而不是从数据库中删除记录
	// 查询时，默认会过滤掉已软删除的记录。如果需要查询已软删除的记录，可以使用 Unscoped().Find() 方法
	// 如果需要永久删除记录，可以使用 Unscoped().Delete() 方法。
}
