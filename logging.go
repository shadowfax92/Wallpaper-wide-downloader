package main

import (
	"io"
	"log"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func LoggingInit(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer, debugHandle io.Writer) {
	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
