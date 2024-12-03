package logHelper

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cast"
)

const (
	ansiReset = "\u001b[0m"

	black   = "0"
	red     = "1"
	green   = "2"
	yellow  = "3"
	blue    = "4"
	magenta = "5"
	cyan    = "6"
	white   = "7"
)

type Tag struct {
	BackgroundColor string
	Color           string
}

type LogFormat struct {
	Level    string
	Message  string
	FuncName string
	Location string
	Time     time.Time
	Tag      Tag
}

var (
	LogTimeZone = "./"
	DebugMode   = false
)

func Warning(message string) {
	timezone, _ := time.LoadLocation(LogTimeZone)
	log := LogFormat{
		Level:    "WARNING",
		Message:  message,
		FuncName: getFuncName(1),
		Location: getPathAndLineNumber(2),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: yellow,
			Color:           black,
		},
	}
	print(log)
}

func Error(message string) {
	timezone, _ := time.LoadLocation(LogTimeZone)
	log := LogFormat{
		Level:    "ERROR",
		Message:  message,
		FuncName: getFuncName(1),
		Location: getPathAndLineNumber(2),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: red,
			Color:           white,
		},
	}
	print(log)
}

func Info(message string) {
	timezone, _ := time.LoadLocation(LogTimeZone)
	log := LogFormat{
		Level:    "INFO",
		Message:  message,
		FuncName: getFuncName(1),
		Location: getPathAndLineNumber(2),
		Time:     time.Now().In(timezone),
		Tag: Tag{
			BackgroundColor: green,
			Color:           black,
		},
	}
	print(log)
}

func print(log LogFormat) {
	logType := createTextBackground(log.Tag.BackgroundColor, log.Tag.Color, " "+log.Level+" ")

	timeNow := log.Time.Format("2006-01-02T15:04:05-0700")
	logTime := createTextBackground(white, black, " "+timeNow+" ")

	logFuncName := createTextBackground(black, white, " "+log.FuncName+" ")

	logLocation := createTextBackground(blue, white, " "+log.Location+" ")

	if DebugMode {
		strLog := logTime + logType + logFuncName + logLocation + " " + log.Message
		fmt.Println(strLog)
	}
}

func createTextBackground(backgroundColor, textColor, text string) string {
	return createBackgroundColor(backgroundColor, false) + createColor(textColor, false) + text + ansiReset
}

func createBackgroundColor(color string, isBright bool) string {
	str := ""
	if isBright {
		str = ";1"
	}

	return "\u001b[4" + color + str + "m"
}

func createColor(color string, isBright bool) string {
	str := ""
	if isBright {
		str = ";1"
	}

	return "\u001b[3" + color + str + "m"
}

func getFuncName(n int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc) // adjust the number
	f := runtime.FuncForPC(pc[0])
	funcName := f.Name()

	arr := strings.Split(funcName, "/")
	sliceArr := arr
	if n <= 0 {
		sliceArr = arr
	} else {
		s := len(arr) - n
		if s < 0 {
			sliceArr = arr
		} else {
			sliceArr = arr[s:]
		}
	}
	strSlice := strings.Join(sliceArr, "/")
	return strSlice
}

func getPathAndLineNumber(n int) string {
	_, filePath, lineNumber, ok := runtime.Caller(2) // adjust the number
	if !ok {
		return "?:?"
	}

	arr := strings.Split(filePath, "/")
	sliceArr := arr
	if n <= 0 {
		sliceArr = arr
	} else {
		s := len(arr) - n
		if s < 0 {
			sliceArr = arr
		} else {
			sliceArr = arr[s:]
		}
	}

	strSlice := strings.Join(sliceArr, "/")
	return strSlice + ":" + cast.ToString(lineNumber)
}
