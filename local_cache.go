package main

import (
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
	"io/ioutil"
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