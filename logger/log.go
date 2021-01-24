package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"time"
)

var Logger *logrus.Logger

func NewLog() *logrus.Logger {
	if Logger != nil {
		return Logger
	}

	writer, err := rotatelogs.New(
		"log-%Y-%m-%d",
		rotatelogs.WithLinkName("log-now"),
		rotatelogs.WithRotationCount(7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}

	Logger = logrus.New()

	Logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%time%] [%lvl%]: %msg% \n",
		},
	))

	Logger.ReportCaller = true
	return Logger
}
