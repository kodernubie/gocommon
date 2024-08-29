package log

// Extends standard log package
// provide automatic configuration using os environment variable
// you can just use golang log, or import this module to use specialize log function
//
// config template :
//
// LOG_LEVEL=fatal|error|info|debug
// 		only accept 1 value
//		error includes fatal
//		info includes fatal + error
//		debug includes fatal + error + error
//
// LOG_FLAG=date|time|microseconds|longfile|shortfile|UTC
//		can accpet multiple value
//		example :
//			log.flag=date|time|shortfile
//
// LOG_OUTPUT=stdout|filepath
//		only accept 1 value
//		example for loging to file
//			log.output=./mylog.log

import (
	"fmt"
	"log"
	"strings"

	"github.com/kodernubie/gocommon/conf"
)

const levelFatal = 0
const levelError = 1
const levelInfo = 2
const levelDebug = 3

var levelLabel []string = []string{"FATAL", "ERROR", "INFO", "DEBUG"}
var level int = levelDebug

func println(targetLevel int, data ...any) {

	if level >= targetLevel {
		data = append([]any{"[" + levelLabel[targetLevel] + "]"}, data...)

		log.Output(3, fmt.Sprintln(data...))
	}
}

func printf(targetLevel int, format string, data ...any) {

	if level >= targetLevel {
		format = "[" + levelLabel[targetLevel] + "] " + format + "\n"
		log.Output(3, fmt.Sprintf(format, data...))
	}
}

func Fatal(data ...any) {
	log.Panic(data...)
}

func Fatalf(format string, data ...any) {
	log.Panicf(format, data...)
}

func Error(data ...any) {
	println(levelError, data...)
}

func Errorf(format string, data ...any) {
	printf(levelError, format, data...)
}

func Info(data ...any) {
	println(levelInfo, data...)
}

func Infof(format string, data ...any) {
	printf(levelInfo, format, data...)
}

func Debug(data ...any) {
	println(levelDebug, data...)
}

func Debugf(format string, data ...any) {
	printf(levelDebug, format, data...)
}

func init() {

	// Level Config
	levelStr := conf.Str("LOG_LEVEL")

	switch levelStr {
	case "fatal":
		level = levelFatal
	case "error":
		level = levelError
	case "info":
		level = levelInfo
	case "debug":
		level = levelDebug
	}

	// Flag Config
	flagStr := conf.Str("LOG_FLAG")

	if flagStr == "" {
		flagStr = "date|time|microseconds|shortfile"
	}

	var flag int

	if strings.Contains(flagStr, "date") {
		flag |= log.Ldate
	}

	if strings.Contains(flagStr, "time") {
		flag |= log.Ltime
	}

	if strings.Contains(flagStr, "microseconds") {
		flag |= log.Lmicroseconds
	}

	if strings.Contains(flagStr, "longfile") {
		flag |= log.Llongfile
	}

	if strings.Contains(flagStr, "shortfile") {
		flag |= log.Lshortfile
	}

	if strings.Contains(flagStr, "UTC") {
		flag |= log.LUTC
	}

	log.SetFlags(flag)
}
