package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"oceanlearn.teach/ginessential/model"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	//host := "192.168.25.110"
	host := viper.GetString("datasource.host")
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	args := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil{
		panic("failed to connect database, err: "+err.Error())
	}

	db.AutoMigrate(&model.User{})  // 自动创建数据表

	DB = db
	return  db
}

func GetDB() *gorm.DB {
	return  DB
}
