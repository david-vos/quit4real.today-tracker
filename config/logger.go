package config

import (
	"log"
	"os"
)

var logger *log.Logger

func Init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func HandleError(message string, err error) {
	if err != nil {
		logger.Printf("%s: %v\n", message, err)
	}
}
