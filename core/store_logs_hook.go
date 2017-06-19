package hoverfly

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type StoreLogsHook struct {
	Entries   []*logrus.Entry
	LogsLimit int
}

func NewStoreLogsHook() *StoreLogsHook {
	return &StoreLogsHook{
		Entries:   []*logrus.Entry{},
		LogsLimit: 1000,
	}
}

func (hook *StoreLogsHook) Fire(entry *logrus.Entry) error {
	if len(hook.Entries) >= hook.LogsLimit {
		hook.Entries = append(hook.Entries[:0], hook.Entries[1:]...)
	}
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

type Fields map[string]interface{}

func (hook StoreLogsHook) GetLogsCount() int {
	return len(hook.Entries)
}

func (hook StoreLogsHook) GetLogs(limit int, from *time.Time) []*logrus.Entry {
	entriesLength := len(hook.Entries)
	if limit > entriesLength {
		limit = entriesLength
	}

	if from != nil {
		entries := []*logrus.Entry{}
		for i := entriesLength - 1; i > 0; i-- {
			if len(entries) < limit {
				entry := hook.Entries[i]
				if entry.Time.After(*from) {
					entries = append([]*logrus.Entry{entry}, entries...)
				}
			}
		}
		return entries
	} else {
		return hook.Entries[entriesLength-limit:]
	}
}
