package db

import (
	"EtuSmartAlarmApi/models"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error while connecting to the database:", err)
		return
	}
	fmt.Println("Connection succeed!")

	DB = DB.Debug() //to be disable during the deployment;

	//Migrations
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Error while migrating Table:", err)
		return
	}
}
