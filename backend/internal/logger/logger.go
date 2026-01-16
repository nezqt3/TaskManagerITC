package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
)

func Init(logFilePath string) error {
	file, err := os.OpenFile(
		logFilePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, file)
	
	Info = log.New(
		mw,
		"[INFO] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Error = log.New(
		mw,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	Fatal = log.New(
		mw,
		"[FATAL] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	return nil
}