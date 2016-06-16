package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
	"path"
)

type HoverflyDirectory struct {
	Path string
}

func NewHoverflyDirectory(config Config) (HoverflyDirectory) {
	if len(config.GetFilepath()) == 0 {
		log.Info("Missing a config file")
		log.Info("Creating a new  a config file")
		hoverflyDirectory := HoverflyDirectory{Path: createHoverflyDirectory(getHomeDirectory())}

		err := config.WriteToFile(hoverflyDirectory)

		if err != nil {
			log.Fatal("Could not write new config to disk")
		}

		return hoverflyDirectory
	}
	return HoverflyDirectory{
		Path: path.Dir(config.GetFilepath()),
	}
}

func getHomeDirectory() (string) {
	homeDirectory, err := homedir.Dir()
	if err != nil {
		log.Debug(err.Error())
		log.Fatal("Unable to get home directory")
	}

	return homeDirectory
}


func createHoverflyDirectory(homeDirectory string)  (string) {
	hoverflyDirectory := filepath.Join(homeDirectory, "/.hoverfly")

	if !fileIsPresent(hoverflyDirectory) {
		err := os.Mkdir(hoverflyDirectory, 0777)
		if err == nil {
			return hoverflyDirectory
		} else {
			log.Debug(err.Error())
			log.Fatal("Could not create a .hoverfly directory")
		}
	}

	return hoverflyDirectory
}