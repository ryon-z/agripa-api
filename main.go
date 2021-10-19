package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "agripa-api/docs" //swagger doccuments

	"agripa-api/common"
	"agripa-api/routers"
)

// .env 파일 로드
func loadEnv() {
	var err error

	if len(os.Args) == 2 {
		err = godotenv.Load(os.Args[1])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal(".env 파일을 불러오지 못했습니다 : " + err.Error())
	}
}

// 초기 설정
func init() {
	loadEnv()

	// DB 연결
	common.DBConnect(os.Getenv("DB_TYPE"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
}

// @title AGRIPA API
// @version 1.0
// @description 아그리파 API
func main() {
	var err error

	//모드 설정
	gin.SetMode(os.Getenv("GIN_MODE"))

	//log 파일 설정
	var logFile *os.File

	logFile, err = os.OpenFile(fmt.Sprintf("log/%s.log", time.Now().Format("20060102")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		os.MkdirAll("log", os.ModePerm)

		logFile, err = os.Create(fmt.Sprintf("log/%s.log", time.Now().Format("20060102")))

		if err != nil {
			fmt.Println(err)
		}
	}

	if gin.Mode() == "release" {
		gin.DefaultWriter = io.MultiWriter(logFile)
	} else {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	}

	r := gin.Default()

	//CORS 설정
	r.Use(cors.Default())

	//Swagger UI 라우팅 설정
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//라우팅 설정
	r = routers.SetRouter(r)

	//서버 On
	endless.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), r)
}
