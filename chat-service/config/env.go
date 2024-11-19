package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
}

var Envs = initConfig()

func initConfig() Config {

	err := godotenv.Load("./chat-service/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", ":8082"),
		DBUser:     getEnv("DB_USER", "tamerlan"),
		DBPassword: getEnv("DB_PASSWORD", "Web_Chat123"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "web_chat"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
