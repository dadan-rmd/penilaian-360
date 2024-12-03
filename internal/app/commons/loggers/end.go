package loggers

import (
	"time"

	"github.com/sirupsen/logrus"
)

// EndRecord for initialize context first time
func EndRecord(record *Data, response string, statuscode int) {
	var level string

	t := time.Since(record.TimeStart)

	if statuscode >= 200 && statuscode < 400 {
		level = "INFO"
	} else if statuscode >= 400 && statuscode < 500 {
		level = "WARN"
	} else {
		level = "ERROR"
	}

	record.StatusCode = statuscode
	record.Response = response
	record.ExecTime = t.Seconds()

	if statuscode == 0 {
		record.StatusCode = 200
	}

	// Getprometheus().MetricRecord(strconv.Itoa(record.StatusCode), record.RequestMethod, record.Endpoint, GetName(), t)

	Output(record, level)
}

// UTCFormatter ...
type UTCFormatter struct {
	logrus.Formatter
}

// Format ...
func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	e.Time = e.Time.In(loc)
	return u.Formatter.Format(e)
}

// Output for output to terminal
func Output(out *Data, level string) {
	logrus.SetFormatter(UTCFormatter{&logrus.JSONFormatter{}})

	if level == "ERROR" {
		logrus.WithField("data", out).Error("apps")
	} else if level == "INFO" {
		logrus.WithField("data", out).Info("apps")
	} else if level == "WARN" {
		logrus.WithField("data", out).Warn("apps")
	}
}
