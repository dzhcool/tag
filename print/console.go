package print

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

var loggerLevel = map[string]int{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"none":  LevelNone,
}

var (
	level int = 3 // 默认级别 error
)

func SetLevel(logLevel string) {
	logLevel = strings.ToLower(logLevel)

	if val, ok := loggerLevel[logLevel]; ok {
		level = val
	}
}

func GetLevel() string {
	for k, v := range loggerLevel {
		if v == level {
			return k
		}
	}
	return ""
}

func br() {
	// fmt.Println("")
}

func Debug(info string) {
	if level > LevelDebug {
		return
	}
	color.Cyan("[debug] " + info)
}

func Debugf(format string, argv ...interface{}) {
	if level > LevelDebug {
		return
	}
	info := fmt.Sprintf(format, argv...)
	Debug(info)
}

func Info(info string) {
	if level > LevelInfo {
		return
	}
	color.Cyan("[info] " + info)
}

func Infof(format string, argv ...interface{}) {
	if level > LevelInfo {
		return
	}
	info := fmt.Sprintf(format, argv...)
	Info(info)
}

func Warn(info string) {
	if level > LevelWarn {
		return
	}
	color.Yellow("[warn] " + info)
}

func Warnf(format string, argv ...interface{}) {
	if level > LevelWarn {
		return
	}
	msg := fmt.Sprintf(format, argv...)
	Warn(msg)
}

func Error(info string) {
	if level > LevelNone {
		return
	}
	color.Red("[error] " + info)
}

func Errorf(format string, argv ...interface{}) {
	if level > LevelNone {
		return
	}
	info := fmt.Sprintf(format, argv...)
	Error(info)
}

func Print(info string) {
	fmt.Printf(info + "\n")
}

func Printf(format string, argv ...interface{}) {
	info := fmt.Sprintf(format, argv...)
	Print(info)
}
