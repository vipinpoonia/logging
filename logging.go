package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

const timeStampFormat = "2006-01-02 15:04:05.999999999Z07:00"

func init() {
	var once sync.Once
	once.Do(func() {
		log.SetFormatter(&log.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: timeStampFormat,
		})
		log.SetLevel(log.TraceLevel)
		log.SetOutput(os.Stdout)
		log.Info("logger is configured")
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
		} else {
			log.AddHook(hook)
		}
	}
	return nil
}
