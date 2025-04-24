package main

import (
	"EtuSmartAlarmApi/db"
	"EtuSmartAlarmApi/services/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	UserStore := user.NewStore(db.DB)

	r := gin.Default()

	user.SetupRoutes(r, user.NewHandler(UserStore))

	port := ":3000" //to be configure in the Env variable
	r.Run(port)
}
