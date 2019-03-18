package configuration

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"strings"

	"net/http"

	log "github.com/sirupsen/logrus"
)

func WriteFile(filePath string, data []byte) error {
	basePath := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	log.Debug(basePath)

	err := os.MkdirAll(basePath, 0644)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(basePath+"/"+fileName, data, 0644)
}

func ReadFile(filePath string) ([]byte, error) {
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return DownloadFile(filePath)
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("File not found: " + filePath)
	}

	return data, nil
}

func DownloadFile(filePath string) ([]byte, error) {
	response, err := http.Get(filePath)
	if err != nil {
		log.Info(err.Error())
		return nil, errors.New("Could not download simulation")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Info(err.Error())
		return nil, errors.New("Could not download simulation")
	}

	return body, nil
}
