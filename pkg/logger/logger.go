package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func NewLogrusLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		PrettyPrint:      false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			funcName := formatFuncName(f.Function)
			fileName := fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
			return funcName, fileName
		},
	})

	return logger
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func formatFuncName(fn string) string {
	parts := strings.Split(fn, ".")
	return parts[len(parts)-1]
}
