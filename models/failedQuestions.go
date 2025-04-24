package models

import (
	"gorm.io/gorm"
)

type FailedQuestion struct {
	gorm.Model
	UserID   *uint  `json:"user_id,omitempty" gorm:"index"`
	Subject  string `json:"subject" gorm:"not null, size: 100" validate:"required, min=2,max=50"`
	Question string `json:"question" gorm:"not null, size: 300" validate:"required, min=10, max=300"`
}
