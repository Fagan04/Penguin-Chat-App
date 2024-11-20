package config

import (
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
		DBUser:     getEnv("DB_USER", "tamerlan"),
		DBPassword: getEnv("DB_PASSWORD", "Web_Chat123"),
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		DBName:     getEnv("DB_NAME", "web_chat"),
		Port:       getEnv("SERVICE_PORT", "8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
