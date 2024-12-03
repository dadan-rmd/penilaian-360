package loggers

import (
	"fmt"
)

// Log ...
func Log(record *Data, message string) {
	record.Messages = append(record.Messages, message)

}

func Logf(record *Data, message string, value ...interface{}) {
	msg := fmt.Sprintf(message, value...)
	record.Messages = append(record.Messages, msg)

}
