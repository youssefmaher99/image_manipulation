package logger

import (
	"log"
	"os"
)

var MyLog *log.Logger

func init() {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	MyLog = l
}
