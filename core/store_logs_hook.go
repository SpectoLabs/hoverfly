package hoverfly

import (
	"github.com/Sirupsen/logrus"
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
	hook.Entries = append([]*logrus.Entry{entry}, hook.Entries...)
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

type Fields map[string]interface{}

func (hook StoreLogsHook) GetLogsCount() int {
	return len(hook.Entries)
}

func (hook StoreLogsHook) GetLogs(limit int) []*logrus.Entry {
	entriesLength := len(hook.Entries)
	if limit > entriesLength {
		limit = entriesLength
	}
	return hook.Entries[:limit]
}
