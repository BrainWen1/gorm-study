// advanced_query/interface.go
package advanced_query

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

func Where_test() {
	var users []UserModel
	err := global.DB.Where("age > ?", 18).Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

func Or_test() {
	var users []UserModel
	err := global.DB.Where("age >= ?", 28).Or("name = ?", "Alice").Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

func Not_test() {
	var users []UserModel
	err := global.DB.Not("age > ?", 18).Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

func Order_test() {
	var users []UserModel
	err := global.DB.Order("age desc").Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

func Scan_test() {
	// Select : 选择特定字段进行查询
	type Result struct {
		Name string
	}

	var results []Result
	// 使用 Select 方法选择 name 字段，并使用 Scan 方法将结果扫描到 results 变量中
	err := global.DB.Model(&UserModel{}).Select("name").Scan(&results).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(results)

	// 也可以直接使用 Scan 方法查询所有字段，并将结果扫描到 results 变量中
	err = global.DB.Model(&UserModel{}).Scan(&results).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(results)
}

func Group_test() {
	type Result struct {
		Age   int
		Count int64
	}

	var results []Result
	err := global.DB.Model(&UserModel{}).Group("age").Select("age, count(*) as count").Scan(&results).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(results)
}

func Distinct_test() {
	type Result struct {
		Age int
	}

	var results []Result
	err := global.DB.Model(&UserModel{}).Distinct("age").Select("age").Scan(&results).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(results)
}

func Page_test(page int, pageSize int) {
	var users []UserModel
	err := global.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}

func Age_18(tx *gorm.DB) *gorm.DB {
	return tx.Where("age = ?", 18)
}

func NameIn(nameList []string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("name IN ?", nameList)
	}
}

func Scope_test() {
	var users []UserModel
	err := global.DB.Scopes(Age_18).Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)

	err = global.DB.Scopes(NameIn([]string{"Alice", "Bob"})).Find(&users).Error
	if err != nil {
		panic("failed to query users")
	}
	fmt.Println(users)
}
