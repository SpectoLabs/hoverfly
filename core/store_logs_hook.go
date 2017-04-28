package hoverfly

import (
	"bytes"

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

func (hook StoreLogsHook) GetLogsView() v2.LogsView {
	return hook.GetFilteredLogsView(500)
}

func (hook StoreLogsHook) GetFilteredLogsView(limit int) v2.LogsView {

	if limit > len(hook.Entries) {
		limit = len(hook.Entries)
	}

	var logs []map[string]interface{}
	for _, entry := range hook.Entries[:limit] {
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

func (hook StoreLogsHook) GetLogs() string {
	return hook.GetFilteredLogs(500)
}

func (hook StoreLogsHook) GetFilteredLogs(limit int) string {

	if limit > len(hook.Entries) {
		limit = len(hook.Entries)
	}

	var buffer bytes.Buffer
	for _, entry := range hook.Entries[:limit] {
		entry.Logger = logrus.New()
		log, err := entry.String()
		if err == nil {
			buffer.WriteString(log)
		}
	}

	return buffer.String()
}
