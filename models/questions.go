package models

import (
	"gorm.io/gorm"
)

type FailedQuestion struct {
	gorm.Model
	UserID   *uint  `json:"user_id,omitempty" gorm:"index"`
	Subject  string `json:"subject" gorm:"not null, size: 100" validate:"required, min=2,max=50"`
	Question string `json:"question" gorm:"not null, size: 300" validate:"required, min=10, max=300"`
	Answers  string `json:"answers" gorm:"not null, size: 500" validate:"required, min=50, max=500"`
}

type Generator_Payload struct {
	UserID            uint   `json:"user_id"`
	Subject           string `json:"subject"`
	NumberOfQuestions uint   `json:"numberofquestions"`
}

type DeepSeekResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
