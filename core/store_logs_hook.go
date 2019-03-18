package hoverfly

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
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
	if hook.LogsLimit == 0 {
		return nil
	}
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

func (hook StoreLogsHook) GetLogs(limit int, from *time.Time) ([]*logrus.Entry, error) {
	if hook.LogsLimit == 0 {
		return []*logrus.Entry{}, fmt.Errorf("Logs disabled")
	}
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
		return entries, nil
	} else {
		return hook.Entries[entriesLength-limit:], nil
	}
}
