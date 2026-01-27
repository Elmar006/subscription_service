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
		log.Println("No .env file found in current directory, trying parent...")
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using system environment")
		}
	}

	cfg := &Config{
		DBHost:     MustGetEnv("DB_HOST"),
		DBPort:     MustGetEnv("DB_PORT"),
		DBUser:     MustGetEnv("DB_USER"),
		DBPass:     MustGetEnv("DB_PASS"),
		DBName:     MustGetEnv("DB_NAME"),
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

func MustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return val
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultValue
}
