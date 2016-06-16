package main

import (
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"os"
	"io/ioutil"
	"errors"
)

type LocalCache struct {
	Uri string
}

func (l *LocalCache) WriteSimulation(simulation Simulation, data []byte) error {
	simulationUri := buildAbsoluteFilePath(l.Uri, simulation.GetFileName())

	return ioutil.WriteFile(simulationUri, data, 0644)
}


func (l *LocalCache) ReadSimulation(simulation Simulation) ([]byte, error) {
	simulationUri := buildAbsoluteFilePath(l.Uri, simulation.GetFileName())

	if !fileIsPresent(simulationUri) {
		return nil, errors.New("Simulation not found")
	}

	return ioutil.ReadFile(simulationUri)
}

func createCacheDirectory(baseUri string) (string, error) {
	cacheDirectory := filepath.Join(baseUri, "cache/")

	if _, err := os.Stat(cacheDirectory); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(cacheDirectory, 0777)
		} else {
			log.Debug(err.Error())
			return "", err
		}
	}

	return cacheDirectory, nil
}

func fileIsPresent(fileUri string) bool {
	if _, err := os.Stat(fileUri); err != nil {
		return os.IsExist(err)
	}

	return true
}

func buildAbsoluteFilePath(baseUri string, fileName string) string {
	return filepath.Join(baseUri, fileName)
}