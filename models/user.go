package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null; unique; size: 50" validate:"required,min=2,max=50"`
	Email    string `json:"email" gorm:"not null; unique; type:varchar(100)" validate:"required,email"`
	Password string `json:"password" gorm:"not null; size:100" validate:"required,min=6"`
	Group    int    `json:"group" validate:"required,gte=2,lte=50000"`
}
