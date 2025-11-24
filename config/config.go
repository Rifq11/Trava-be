package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	// dsn := "root@tcp(host.docker.internal:3306)/trava?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root@tcp(127.0.0.1:3306)/trava?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		panic(err.Error())
	}

	DB = db
}