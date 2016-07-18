package main

import (
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"errors"
)

type LogFile struct {
	Path string
	Name string
}

func NewLogFile(directory HoverflyDirectory, adminPort, proxyPort string) (LogFile) {
	fileName := fmt.Sprintf("hoverfly.%v.%v.log", adminPort, proxyPort)

	filePath := filepath.Join(directory.Path, fileName)

	return LogFile{
		Path: filePath,
		Name: fileName,
	}
}

func (l *LogFile) GetLogs() (string, error) {
	content, err := ioutil.ReadFile(l.Path)
	if err != nil {
		log.Debug(err.Error())
		return "", errors.New("Could not open Hoverfly log file")
	}

	return string(content), nil
}