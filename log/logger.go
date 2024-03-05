/*
Just a plain logger, its
*/
package log

import (
	"bhd/dailyrotate"
	"fmt"
	"strings"
	"time"
)

const (
	LevelSql     = -1
	LevelFine    = 0
	LevelDebug   = 1
	LevelInfo    = 2
	LevelWarning = 3
	LevelError   = 4
)

type Log struct {
	dest      *dailyrotate.File
	logPrefix string
	logLevel  int
}

var l Log = Log{logPrefix: "BHD v1.0:"}

func SetLogLevel(level int) {
	if level < 0 {
		level = LevelSql
	}
	if level >= 4 {
		level = LevelError
	}
	l.logLevel = level
}

func GetLogLevel() int {
	return l.logLevel
}

func SetLogDest(file *dailyrotate.File) {
	l.dest = file
}

func SetLogPrefix(prefix string) {
	l.logPrefix = prefix
}

func log(lvl int, a ...interface{}) {
	var ll string
	switch lvl {
	case -1:
		ll = "[SQL]"
	case 0:
		ll = "[FIN]"
	case 1:
		ll = "[DEB]"
	case 2:
		ll = "[INF]"
	case 3:
		ll = "[WARN]"
	case 4:
		ll = "[ERR]"
	}
	// date time
	st := fmt.Sprint(a...)
	st = strings.TrimRight(strings.TrimLeft(st, "["), "]")
	tx := time.Now().Format(time.RFC3339)
	tx = strings.Replace(tx, "T", " ", 1)
	fmt.Println(tx[:len(tx)-2], ll, l.logPrefix, st)
	if l.dest != nil {
		fmt.Fprintln(l.dest, tx[:len(tx)-2], ll, l.logPrefix, st)
	}
}

func Sql(a ...interface{}) {
	if l.logLevel <= LevelSql && l.dest != nil {
		log(LevelSql, a)
	}
}

func Fine(a ...interface{}) {
	if l.logLevel <= LevelFine && l.dest != nil {
		log(LevelFine, a)
	}
}

func Debug(a ...interface{}) {
	if l.logLevel <= LevelDebug && l.dest != nil {
		log(LevelDebug, a)
	}
}

func Info(a ...interface{}) {
	if l.logLevel <= LevelInfo && l.dest != nil {
		log(LevelInfo, a)
	}
}

func Warn(a ...interface{}) {
	if l.logLevel <= LevelWarning && l.dest != nil {
		log(LevelWarning, a)
	}
}

func Error(a ...interface{}) {
	if l.logLevel <= LevelError && l.dest != nil {
		var stackTrace = ""
		for _, el := range a {
			if _, ok := el.(error); ok {
				stackTrace = fmt.Sprintf("%+v", el)
				a = append(a, stackTrace)
			}
		}
		log(LevelError, a)
	}
}
