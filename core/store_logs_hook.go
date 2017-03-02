package hoverfly

import (
	"github.com/Sirupsen/logrus"
	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

type StoreLogsHook struct {
	Entries []*logrus.Entry
}

func NewStoreLogsHook() *StoreLogsHook {
	return &StoreLogsHook{
		Entries: []*logrus.Entry{},
	}
}

func (hook *StoreLogsHook) Fire(entry *logrus.Entry) error {
	hook.Entries = append(hook.Entries, entry)
	return nil
}

func (hook StoreLogsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook StoreLogsHook) GetLogsView() v2.LogsView {
	var logs []v2.LogView
	for _, entry := range hook.Entries {
		logs = append(logs, v2.LogView{
			Time:    entry.Time.Format(logrus.DefaultTimestampFormat),
			Message: entry.Message,
			Level:   entry.Level.String(),
		})
	}

	return v2.LogsView{
		Logs: logs,
	}
}
