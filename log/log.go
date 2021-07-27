package log

import (
	golog "log"
	"os"
)

var loggers = []Logger{golog.New(os.Stderr, "", golog.LstdFlags)}

func UseLogger(l Logger) {
	loggers = append([]Logger{l}, loggers...)
}

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

func Print(v ...interface{}) {
	for _, l := range loggers {
		l.Print(v...)
	}
}

func Printf(format string, v ...interface{}) {
	for _, l := range loggers {
		l.Printf(format, v...)
	}
}

func Println(v ...interface{}) {
	for _, l := range loggers {
		l.Println(v...)
	}
}

func Fatal(v ...interface{}) {
	for _, l := range loggers {
		l.Fatal(v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	for _, l := range loggers {
		l.Fatalf(format, v...)
	}
}

func Fatalln(v ...interface{}) {
	for _, l := range loggers {
		l.Fatalln(v...)
	}
}

func Panic(v ...interface{}) {
	for _, l := range loggers {
		l.Panic(v...)
	}
}

func Panicf(format string, v ...interface{}) {
	for _, l := range loggers {
		l.Panicf(format, v...)
	}
}

func Panicln(v ...interface{}) {
	for _, l := range loggers {
		l.Panicln(v...)
	}
}
