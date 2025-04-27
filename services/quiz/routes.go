package quiz

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, handler *Handler) {
	public := r.Group("/my_alarm/v1")
	{
		public.POST("/quiz/submit_wrong_answers", handler.saving)
		public.POST("/quiz/generate_quiz", handler.QuizGenerator)
	}
}
