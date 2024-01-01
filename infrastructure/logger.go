package infrastructure

import (
	"go-learning/clean-architecture-mongo/usecases"
	"io"
	"log"
	"os"
)

type Logger struct{}

func NewLogger() usecases.Logger {
	return &Logger{}
}

func (l Logger) LogError(format string, v ...interface{}) {
	file, err := os.OpenFile("./log/error.log", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("%s", err)
	}

	defer file.Close()
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf(format, v...)
}

func (l Logger) LogAccess(format string, v ...interface{}) {
	file, err := os.OpenFile("./log/access.log", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("%s", err)
	}

	defer file.Close()
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)
	log.Printf(format, v...)
}
