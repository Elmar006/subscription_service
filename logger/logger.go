package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetOutput(os.Stdout)

	if os.Getenv("APP_ENV") == "prod" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006.01.02 15:04:05",
		})
	}

	log.SetLevel(log.InfoLevel)
	if os.Getenv("APP_ENV") == "dev" {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}
}

func L() *log.Logger {
	return log.StandardLogger()
}
