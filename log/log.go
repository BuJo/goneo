package log

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[goneo] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.LUTC|log.Lmsgprefix)
}

func SetDefault(l *log.Logger) {
	logger = l
}

func Fatal(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Print(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
}

func Println(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
}

func Printf(format string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, v...))
}
