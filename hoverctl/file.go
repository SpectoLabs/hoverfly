package main

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
)

func WriteFile(filePath string, data []byte) error {
	basePath := path.Dir(filePath)
	fileName := path.Base(filePath)
	log.Debug(basePath)

	err := os.MkdirAll(basePath, 0644)
	if err != nil {
		return err
	}
	log.Debug("Should be writing?")
	return ioutil.WriteFile(basePath+"/"+fileName, data, 0644)
}

func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}
