package main

import (
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"errors"
)

type LocalCache struct {
	uri string
}

func (l *LocalCache) PersistSimulation(key string, data []byte) error {
	vendor, name := splitHoverfileName(key)

	hoverfileName := buildHoverfileName(vendor, name)

	hoverfileUri := buildHoverfileUri(hoverfileName, l.uri)

	return ioutil.WriteFile(hoverfileUri, data, 0644)
}


func (l *LocalCache) ReadSimulation(key string) ([]byte, error) {
	vendor, name := splitHoverfileName(key)

	hoverfileName := buildHoverfileName(vendor, name)
	hoverfileUri := buildHoverfileUri(hoverfileName, l.uri)

	if !fileIsPresent(hoverfileUri) {
		return nil, errors.New("Simulation not found")
	}

	return ioutil.ReadFile(hoverfileUri)
}

func createHomeDirectory() (string, error) {
	homeDirectory, _ := homedir.Dir()
	hoverflyDirectory := filepath.Join(homeDirectory, "/.hoverfly")


	if !fileIsPresent(hoverflyDirectory) {
		err := os.Mkdir(hoverflyDirectory, 0777)
		if err == nil {
			return hoverflyDirectory, err
		} else {
			return "", err
		}
	}

	return hoverflyDirectory, nil
}

func createCacheDirectory(baseUri string) (string, error) {
	cacheDirectory := filepath.Join(baseUri, "cache/")

	if _, err := os.Stat(cacheDirectory); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(cacheDirectory, 0777)
		} else {
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

func buildHoverfileName(vendor string, api string) string {
	return fmt.Sprintf("%v.%v.hfile", vendor, api)
}

func splitHoverfileName(hoverfileKey string) (string, string) {
	s := strings.Split(hoverfileKey, "/", )
	vendor := s[0]
	name := s[1]

	return vendor, name
}