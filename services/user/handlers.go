package user

import (
	"EtuSmartAlarmApi/models"
	"EtuSmartAlarmApi/services/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

type Handler struct {
	Store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{Store: store}
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

var validate = validator.New()

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		// On envoie l'erreur détaillée au client pour debug
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Data validation avec messages clairs
	if err := validate.Struct(newUser); err != nil {
		errors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(errors))
		for i, e := range errors {
			errorMessages[i] = fmt.Sprintf("Field %s failed validation: %s", e.Field(), e.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server problem"})
		fmt.Println("Hashing error:", err)
		return
	}
	newUser.Password = hashedPassword

	// Enregistrement en base
	if err := h.Store.db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created!"})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func (h *Handler) LoginUser(c *gin.Context) {
	var LoginUserPayload models.LoginUserPayloads
	if err := c.ShouldBindJSON(&LoginUserPayload); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "Invalid input"})
		return
	}

	if err := validate.Struct(LoginUserPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		errorMessages := make([]string, len(errors))
		for i, e := range errors {
			errorMessages[i] = fmt.Sprintf("Field %s failed validation: %s", e.Field(), e.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
		return
	}

	var user models.User
	if err := h.Store.db.Where("username = ?", LoginUserPayload.Username).First(&user).Error; err != nil {
		log.Printf("Failed to find user with email %s: %v", LoginUserPayload.Username, err)
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Email or Password"})
		return
	}

	ok := auth.ComparePassword(user.Password, []byte(LoginUserPayload.Password))
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid Email or Password"})
		return
	}

	//retourner le username et le numero de groupe
	user.Password = "********"

	c.JSON(200, user)
}
