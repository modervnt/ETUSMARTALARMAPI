package models

type LoginUserPayloads struct {
	//ca aurait ete mieux de send le username
	Username string `json:"username" gorm:"not null; unique; size: 50" validate:"required,min=2,max=50"`
	Password string `json:"password" gorm:"not null; size:100" validate:"required,min=6"`
}
