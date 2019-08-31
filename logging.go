package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"os"
	"path"
	"runtime"
	"sync"
)

const timeStampFormat = "2006-01-02 15:04:05.999Z07:00"

func init() {
	var once sync.Once
	once.Do(func() {

		log.SetLevel(log.TraceLevel)
		log.SetOutput(os.Stdout)
		log.Info("logger is configured")
		log.SetReportCaller(true)

		log.SetFormatter(&log.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: timeStampFormat,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return "", fmt.Sprintf(" %s:%d", filename, f.Line)
			},
		})
	})
}

func Init(conf LogConfig) error {
	if conf.SentryDsn == "" {
		log.Warn("SentryDsn is missing")
	} else {
		sentryOptions := Options{Dsn: conf.SentryDsn, AttachStacktrace: true, Environment: conf.Env}
		hook, err := NewHook(sentryOptions, log.PanicLevel, log.FatalLevel, log.ErrorLevel)
		if err != nil {
			log.Error("failed to add hook for sentry err:", err)
			return errors.New(err.Error())
		} else {
			log.AddHook(hook)
		}
	}
	return nil
}
