// reference code https://github.com/onrik/logrus.git

package logging

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var (
	levelsMap = map[logrus.Level]sentry.Level{
		logrus.PanicLevel: sentry.LevelFatal,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.TraceLevel: sentry.LevelDebug,
	}
)

type Hook struct {
	client      *sentry.Client
	levels      []logrus.Level
	tags        map[string]string
	release     string
	environment string
}

type Options sentry.ClientOptions

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) SetTags(tags map[string]string) {
	hook.tags = tags
}

func (hook *Hook) AddTag(key, value string) {
	hook.tags[key] = value
}

func (hook *Hook) SetRelease(release string) {
	hook.release = release
}

func (hook *Hook) SetEnvironment(environment string) {
	hook.environment = environment
}

func (hook *Hook) constructEvent(entry *logrus.Entry) *sentry.Event{
	var exceptions []sentry.Exception
	err, ok := entry.Data[logrus.ErrorKey].(error)

	event := sentry.Event{
		Level:       levelsMap[entry.Level],
		Extra:       map[string]interface{}(entry.Data),
		Tags:        hook.tags,
		Environment: hook.environment,
		Release:     hook.release,
		Exception:   exceptions,
		Message: entry.Message,
	}
	if !ok {
		event.Threads = []sentry.Thread{{
			Stacktrace: sentry.NewStacktrace(),
			Crashed:    false,
			Current:    true,
		}}
		return &event
	}
	stacktrace := sentry.ExtractStacktrace(err)
	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}
	exceptions = append(exceptions, sentry.Exception{
		Type:       entry.Message,
		Value:      err.Error(),
		Stacktrace: stacktrace,
	})

	event.Exception = exceptions
	return &event

}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	fmt.Println(entry.Caller.File, entry.Caller.Line)
	event := hook.constructEvent(entry)

	hub := sentry.CurrentHub()
	hook.client.CaptureEvent(event, nil, hub.Scope())
	return nil
}

func NewHook(options Options, levels ...logrus.Level) (*Hook, error) {
	client, err := sentry.NewClient(sentry.ClientOptions(options))
	if err != nil {
		return nil, err
	}
	hook := Hook{
		client: client,
		levels: levels,
		tags:   map[string]string{},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook, nil
}
