package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动，不直接使用，但必须导入以注册驱动
)

type User struct {
	ID    int
	Name  string
	age   int
	email string
	phone string
	addr  string
}

func main() {
	// 数据源名称 (DSN)，格式为 "username:password@protocol(address)/dbname"
	dsn := "root:123456@tcp(192.168.12.143:3306)/game"

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// 操作数据库
	_, err = db.Exec("use game") // 选择数据库
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("select name,id from users") // 执行查询
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&name, &id); err != nil { // 扫描结果到变量
			panic(err)
		}
		fmt.Println(name, id) // 输出查询结果
	}

	var name string
	err = db.QueryRow("select name from users where age=18").Scan(&name) // 查询单行结果，如果有多行匹配，只返回表中第一行的结果
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s\n", name) // 输出查询结果
}
