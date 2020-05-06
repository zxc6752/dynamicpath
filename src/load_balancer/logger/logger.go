package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

var log *logrus.Logger
var InitLog *logrus.Entry
var ContextLog *logrus.Entry
var UtilLog *logrus.Entry
var MonitorClientLog *logrus.Entry
var SmfClientLog *logrus.Entry
var HandlerLog *logrus.Entry

func init() {
	log = logrus.New()

	log.Formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			orgFilename, _ := os.Getwd()
			repopath := orgFilename
			repopath = strings.Replace(repopath, "/bin", "", 1)
			filename := strings.Replace(f.File, repopath, "", -1)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	InitLog = log.WithFields(logrus.Fields{"loadBalancer": "init"})
	ContextLog = log.WithFields(logrus.Fields{"loadBalancer": "context"})
	UtilLog = log.WithFields(logrus.Fields{"loadBalancer": "util"})
	MonitorClientLog = log.WithFields(logrus.Fields{"loadBalancer": "monitor"})
	HandlerLog = log.WithFields(logrus.Fields{"loadBalancer": "handler"})
	SmfClientLog = log.WithFields(logrus.Fields{"loadBalancer": "handler"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
