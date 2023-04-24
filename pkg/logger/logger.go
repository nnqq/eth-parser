package logger

import "log"

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
		log.Printf("%s: "+format+"\n", l.pkg, v)
	}
}
