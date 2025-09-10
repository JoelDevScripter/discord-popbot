package logger

import (
	"io"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	// Abrir archivo
	file, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// MultiWriter para consola + archivo
	mwInfo := io.MultiWriter(os.Stdout, file)
	mwError := io.MultiWriter(os.Stderr, file)

	infoLogger = log.New(mwInfo, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(mwError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

func Error(v ...interface{}) {
	errorLogger.Println(v...)
}
