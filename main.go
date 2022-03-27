package main

import (
	"fmt"
	"os"
	"quik/config"
	"quik/models"
	"quik/routers"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

//func init runs before the main function
func init() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Println(".env file missing")
	}
	config.DB, err = gorm.Open(mysql.Open(config.DbURL(config.BuildDbConfig())), &gorm.Config{})
	if err != nil {
		fmt.Println("status: ", err)
	}

	config.DB.AutoMigrate(&models.User{}, &models.UserWallet{})

	formatter := runtime.Formatter{ChildFormatter: &log.JSONFormatter{}}
	formatter.Line = true
	log.SetFormatter(&formatter)
	log.SetOutput(os.Stdout)

}
func main() {

	r := routers.SetUpRouter()
	r.Run()
}
