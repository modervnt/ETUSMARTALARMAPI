package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost       string
	Port             string
	DBUser           string
	DBPassword       string
	DBAddress        string
	DBName           string
	DEEPSEEK_API_KEY string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:       getEnv("PUBLIC_HOST", "http://localhost"),
		Port:             getEnv("PORT", "3000"),
		DBUser:           getEnv("DB_USER", "root"),
		DBPassword:       getEnv("DB_PASSWORD", "mypassword"),
		DBAddress:        fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:           getEnv("DB_NAME", "data.db"),
		DEEPSEEK_API_KEY: getEnv("DEEPSEEK_API_KEY", "sk-1482244387914112b36dae7fa58d1382"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
