package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	ServerPort string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("env file not found, using system enviroment")
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", "postgres"),
		DBName:     getEnv("DB_NAME", "subscriptions"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	return cfg
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPass,
		c.DBName,
	)
}

func MusGetEnv(key, value string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("Required enviroment variable %s is not set", key)
	}

	return val
}

func getEnv(key, value string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return value
}
