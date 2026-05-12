package main

import (
	"gorm.io/gorm"

	"gorm-study/06-Transaction/global"
)

type User struct {
	gorm.Model // gorm.Model 包含 ID、CreatedAt、UpdatedAt 和 DeletedAt 字段
	Name       string
}

type UserDetail struct {
	ID      int
	Address string
	UserID  uint
}

func migrate() {
	err := global.DB.AutoMigrate(&User{}, &UserDetail{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func CreateUserWithDetail(name string, address string) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		user := User{Name: name}
		err := tx.Create(&user).Error
		if err != nil {
			return err
		}

		// 创建用户详情
		detail := UserDetail{Address: address, UserID: user.ID}
		err = tx.Create(&detail).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func CreateUserWithDetailManual(name string, address string) error {
	// 获取数据库连接
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 注册recover函数，确保在发生panic时能够回滚事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	user := User{Name: name}
	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 创建用户详情
	detail := UserDetail{Address: address, UserID: user.ID}
	err = tx.Create(&detail).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func main() {
	global.ConnectDB()
	migrate()

	// 自动事务
	err := CreateUserWithDetail("Xiao Ming", "Beijing")
	if err != nil {
		panic("failed to create user with detail")
	}

	// 手动事务
	err = CreateUserWithDetailManual("Da Hong", "Shanghai")
	if err != nil {
		panic("failed to create user with detail")
	}
}
