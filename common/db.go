package common

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver import
)

// DB 객체
var DB *gorm.DB

// DBConnect 데이터베이스 연결
func DBConnect(dbType string, dbHost string, dbName string, dbUser string, dbPassword string) *gorm.DB {
	db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName))

	if err != nil {
		fmt.Println("DB connection Error: ", err)
	}

	DB = db

	return DB
}

// DBDisConnect 데이터베이스 연결해제
func DBDisConnect() error {
	err := DB.Close()
	return err
}

// GetDB 데이터베이트 객체 조회
func GetDB() *gorm.DB {
	return DB
}
