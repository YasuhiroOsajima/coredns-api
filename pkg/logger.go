package pkg

import (
	"log"
	"os"
)

var logfile = "/tmp/test.log"

func SetLogger(processUUID string) *log.Logger {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	logger := log.New(f, "["+processUUID+"]", log.LstdFlags|log.LUTC)
	return logger
}
