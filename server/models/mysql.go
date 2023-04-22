package models

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ChatDB *gorm.DB

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := viper.GetString(`mysql.dsn`)
	ChatDB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	AutoCreateTable()
	return
}

func AutoCreateTable() {
	err := ChatDB.AutoMigrate(&Message{})
	if err != nil {
		log.Println(err)
	}
	err = ChatDB.AutoMigrate(&User{})
	if err != nil {
		log.Println(err)
	}
}
