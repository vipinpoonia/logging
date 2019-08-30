package logging

import (
	log "github.com/sirupsen/logrus"
	"logging/sentry"
	"os"
)

const timeStampFormat = "2006-01-02 15:04:05.999999999Z07:00"

func Init(conf LogConfig) error {

	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: timeStampFormat,
	})
	log.SetLevel(log.TraceLevel)
	log.SetOutput(os.Stdout)
	log.Info("logger is configured")

	if conf.SentryDsn == "" {
		log.Warn("SentryDsn is missing")
	} else {
		sentryOptions := sentry.Options{Dsn: conf.SentryDsn, AttachStacktrace: true, Environment: conf.Env}
		hook, err := sentry.NewHook(sentryOptions, log.PanicLevel, log.FatalLevel, log.ErrorLevel)
		if err != nil {
			log.Error("failed to add hook for sentry err:", err)
		} else {
			log.AddHook(hook)
		}
	}
	return nil
}
