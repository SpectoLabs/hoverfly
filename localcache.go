package main

import (
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"os"
	"io/ioutil"
	"errors"
)

type LocalCache struct {
	URI string
}

func (l *LocalCache) WriteSimulation(simulation Simulation, data []byte) error {
	simulationURI := buildAbsoluteFilePath(l.URI, simulation.GetFileName())

	return ioutil.WriteFile(simulationURI, data, 0644)
}


func (l *LocalCache) ReadSimulation(simulation Simulation) ([]byte, error) {
	simulationURI := buildAbsoluteFilePath(l.URI, simulation.GetFileName())

	if !fileIsPresent(simulationURI) {
		return nil, errors.New("Simulation not found")
	}

	return ioutil.ReadFile(simulationURI)
}

func createCacheDirectory(hoverflyDirectory HoverflyDirectory) (string, error) {
	cacheDirectory := filepath.Join(hoverflyDirectory.Path, "cache/")

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

func fileIsPresent(fileURI string) bool {
	if _, err := os.Stat(fileURI); err != nil {
		return os.IsExist(err)
	}

	return true
}

func buildAbsoluteFilePath(baseURI string, fileName string) string {
	return filepath.Join(baseURI, fileName)
}