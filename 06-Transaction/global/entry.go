// global/entry.go
package global

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:123456@tcp(192.168.12.143:3306)/multi_tab"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
