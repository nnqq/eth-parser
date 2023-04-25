package logger

import (
	"fmt"
	"log"
)

type Printer interface {
	Printf(format string, v ...interface{})
}

type Logger struct {
	pkg     string
	enabled bool
}

func NewLogger(pkg string, enabled bool) *Logger {
	return &Logger{
		pkg:     pkg,
		enabled: enabled,
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.enabled {
		log.Println(l.pkg + ": " + fmt.Sprintf(format, v...))
	}
}
