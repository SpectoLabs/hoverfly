package wrapper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
)

type LogFile struct {
	Path string
	Name string
}

func NewLogFile(directory HoverflyDirectory, adminPort, proxyPort string) LogFile {
	fileName := fmt.Sprintf("hoverfly.%v.%v.log", adminPort, proxyPort)

	filePath := filepath.Join(directory.Path, fileName)

	return LogFile{
		Path: filePath,
		Name: fileName,
	}
}

func (l *LogFile) getLogs() (string, error) {
	content, err := ioutil.ReadFile(l.Path)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not open Hoverfly log file")
	}

	return string(content), nil
}

func (l *LogFile) Print() error {
	logs, err := l.getLogs()

	if err != nil {
		return err
	}

	fmt.Print(logs)

	return nil
}

func (l *LogFile) Tail() error {
	tail, err := tail.TailFile(l.Path, tail.Config{Follow: true})
	if err != nil {
		log.Debug(err.Error())
		return errors.New("Could not follow Hoverfly log file")
	}

	for line := range tail.Lines {
		fmt.Println(line.Text)
	}

	return nil
}
