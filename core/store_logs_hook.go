package hoverfly

import (
	"github.com/Sirupsen/logrus"
	"time"
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
					entries = append(entries, entry)
				}
			}
		}
		return entries
	} else {
		return hook.Entries[entriesLength-limit:]
	}
}
