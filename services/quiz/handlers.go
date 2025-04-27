package quiz

import (
	"EtuSmartAlarmApi/configs"
	"EtuSmartAlarmApi/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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

func callDeepSeekAPI(prompt string) (string, error) {
	client := resty.New()

	//apiKey := "Bearer " + configs.Envs.DEEPSEEK_API_KEY
	apiKey := fmt.Sprintf("Bearer %s", configs.Envs.DEEPSEEK_API_KEY)
	apiUrl := "https://api.deepseek.com/v1/chat/completions" // Notez le https://

	// Configuration du timeout
	client.SetTimeout(30 * time.Second)

	resp, err := client.R().
		SetHeader("Authorization", apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "deepseek-chat",
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"temperature": 0.7,  // Contrôle la créativité
			"max_tokens":  1000, // Limite la longueur de la réponse
		}).
		Post(apiUrl)

	if err != nil {
		return "", fmt.Errorf("échec de la requête API: %v", err)
	}

	// Journalisation du corps de la réponse pour le débogage
	log.Printf("Réponse API - Status: %d, Body: %s", resp.StatusCode(), resp.Body())

	if resp.StatusCode() != http.StatusOK {
		var apiErr struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}

		if err := json.Unmarshal(resp.Body(), &apiErr); err == nil && apiErr.Error.Message != "" {
			return "", fmt.Errorf("erreur API (%d): %s - Type: %s",
				resp.StatusCode(), apiErr.Error.Message, apiErr.Error.Type)
		}

		return "", fmt.Errorf("erreur API (%d): %s", resp.StatusCode(), resp.Body())
	}

	return string(resp.Body()), nil
}

func ParseQuestions(jsonResponse string) string {
	var response models.DeepSeekResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		return content
	} else {
		return ""
	}

}

func (h *Handler) QuizGenerator(c *gin.Context) {
	var generatorPayload models.Generator_Payload
	if err := c.ShouldBindJSON(&generatorPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Récupérer les questions échouées
	var failedQuestions []models.FailedQuestion
	if err := h.Store.db.
		Where("user_id = ? AND subject = ?", generatorPayload.UserID, generatorPayload.Subject).
		Limit(int(generatorPayload.NumberOfQuestions)).
		Find(&failedQuestions).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to get failed questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	// Calculer le nombre de questions à générer
	toBeGenerated := int(generatorPayload.NumberOfQuestions) - len(failedQuestions)
	if toBeGenerated < 0 {
		toBeGenerated = 0
	}

	var generatedQuestions string
	if toBeGenerated > 0 {
		prompt := fmt.Sprintf(`Generate %d MCQs about "%s" using this exact format:
Question [X]: [Text]
A) [Option1]
B) [Option2]
C) [Option3]
D) [Option4]
Correct: [Letter]

Rules:
- No Markdown/LaTeX
- 4 options per question
- 1 correct answer
- Blank line between questions`, toBeGenerated, generatorPayload.Subject)

		apiResponse, err := callDeepSeekAPI(prompt)
		if err != nil {
			log.Printf("API call failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate questions"})
			return
		}

		generatedQuestions = ParseQuestions(apiResponse)
		if generatedQuestions == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse generated questions"})
			return
		}
	}

	// Combiner les questions
	var fullQuiz string
	if generatedQuestions != "" {
		fullQuiz = generatedQuestions
	}

	// Ajouter les questions échouées
	for _, q := range failedQuestions {
		fullQuiz += fmt.Sprintf("\n\n%s", q.Question)
	}

	c.JSON(http.StatusOK, gin.H{"quiz": fullQuiz})
}

func (h *Handler) saving(c *gin.Context) {
	var failedQuestions []models.FailedQuestion
	if err := c.ShouldBindJSON(&failedQuestions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.Store.db.Create(&failedQuestions).Error; err != nil {
		log.Printf("Failed to save questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save questions"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Questions saved successfully"})
}
