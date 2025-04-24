package user

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, handler *Handler) {
	public := r.Group("/my_alarm/v1")
	{
		public.POST("/users/register", handler.CreateUser)
		public.POST("/users/login", handler.LoginUser)
	}
}
