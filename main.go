package main

import (
	"EtuSmartAlarmApi/configs"
	"EtuSmartAlarmApi/db"
	"EtuSmartAlarmApi/services/quiz"
	"EtuSmartAlarmApi/services/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	UserStore := user.NewStore(db.DB)
	QuizStore := quiz.NewStore(db.DB)

	r := gin.Default()

	r.Use(cors.Default())

	user.SetupRoutes(r, user.NewHandler(UserStore))
	quiz.SetupRoutes(r, quiz.NewHandler(QuizStore))

	port := ":" + configs.Envs.Port
	r.Run(port)
}
