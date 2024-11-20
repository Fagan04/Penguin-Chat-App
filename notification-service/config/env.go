package config

import (
	"fmt"
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

var Envs = LoadConfig()

func LoadConfig() *Config {
	return &Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		DBUser:     getEnv("DB_USER", "tamerlan"),
		DBPassword: getEnv("DB_PASSWORD", "Web_Chat123"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "web_chat"),
		Port:       getEnv("PORT", "8082"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
