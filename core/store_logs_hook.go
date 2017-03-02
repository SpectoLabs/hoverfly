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

type Fields map[string]interface{}

func (hook StoreLogsHook) GetLogsView() v2.LogsView {

	var logs []map[string]interface{}
	for _, entry := range hook.Entries {
		data := make(map[string]interface{}, len(entry.Data)+3)

		for k, v := range entry.Data {
			data[k] = v
		}

		data["time"] = entry.Time.Format(logrus.DefaultTimestampFormat)
		data["msg"] = entry.Message
		data["level"] = entry.Level.String()

		logs = append(logs, data)
	}

	return v2.LogsView{
		Logs: logs,
	}
}
