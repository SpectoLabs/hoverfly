package util

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type LogFileConfig struct {
	Filename  string
	Level     logrus.Level
	Formatter logrus.Formatter
}

type LogFileHook struct {
	Config    LogFileConfig
	logWriter io.Writer
}

func NewLogFileHook(config LogFileConfig) (logrus.Hook, error) {

	hook := LogFileHook{
		Config: config,
	}
	file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	hook.logWriter = file

	return &hook, nil
}

func (hook *LogFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *LogFileHook) Fire(entry *logrus.Entry) (err error) {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.logWriter.Write(b)
	return nil
}
