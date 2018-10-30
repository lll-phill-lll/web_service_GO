package logger

import (
	"log"
	"io"
	"os"
)

// https://www.ardanlabs.com/blog/2013/11/using-log-package-in-go.html?showComment=1396035887595

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func SetLogger(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	file, err := os.OpenFile("FlyLogs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)	// file for logs
	if err != nil {
		log.Println("Failed to open log file:", err)
	}

	outTrace := io.MultiWriter(file, traceHandle)	// copy logs stream to file
	outInfo:= io.MultiWriter(file, infoHandle)
	outWarning := io.MultiWriter(file, warningHandle)
	outError := io.MultiWriter(file, errorHandle)

	// set function to write log, example: Info.Println("someText")
	Trace = log.New(outTrace,
		"TRACE: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Info = log.New(outInfo,
		"INFO: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Warning = log.New(outWarning,
		"WARNING: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Error = log.New(outError,
		"ERROR: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
}
