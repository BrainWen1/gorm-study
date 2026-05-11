package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm-study/05-Custom-datatype/global"
)

// 枚举类型
type Status int

const (
	StatusActive   Status = iota // 活跃状态
	StatusInactive               // 非活跃状态
)

func (s *Status) Scan(value interface{}) error {
	*s = Status(value.(int64))
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return int64(s), nil
}

// json结构体
type Info struct {
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// User 结构体，包含一个 Info 字段，使用自定义数据类型存储为 JSON 格式
type User struct {
	ID     int
	Name   string
	Info   Info   `gorm:"type:json"` // 自定义数据类型，存储为 JSON 格式
	Status Status `gorm:"type:int"`  // 使用枚举类型存储状态
}

// Scan 实现 sql.Scanner 接口，用于从数据库中扫描值
func (u *Info) Scan(value interface{}) error {
	byte, ok := value.([]byte) // 将数据库中的值转换为 []byte 类型，即字节数组
	if !ok {
		return fmt.Errorf("failed to convert value to []byte")
	}

	err := json.Unmarshal(byte, u) // 将字节数组中的 JSON 数据解析到 Info 结构体中
	if err != nil {
		return err
	}
	return nil
}

func (u Info) Value() (driver.Value, error) {
	byte, err := json.Marshal(u) // 将 Info 结构体转换为 JSON 格式的字节数组
	if err != nil {
		return nil, err
	}
	return byte, nil
}

// MarshalJSON 实现 json.Marshaler 接口，用于自定义 JSON 序列化
func (s Status) MarshalJSON() (data []byte, err error) {
	// 根据枚举值设置对应的字符串描述
	var str string
	switch s {
	case StatusActive:
		str = "活跃"
	case StatusInactive:
		str = "非活跃"
	}
	// 将枚举值和对应的字符串一起序列化为 JSON 格式
	return json.Marshal(map[string]any{
		"status":      int(s),
		"statusTitle": str,
	})
}

func main() {
	global.ConnectDB()
	global.DB.AutoMigrate(&User{})

	// 插入数据
	err := global.DB.Create(&User{
		Name: "Alice",
		Info: Info{
			Email: "alice@gmail.com",
			Age:   18,
		},
		Status: StatusActive,
	}).Error
	if err != nil {
		panic("failed to insert data")
	}

	// 查询数据
	var user User
	err = global.DB.First(&user).Error
	if err != nil {
		panic("failed to query data")
	}
	fmt.Println(user)
}
