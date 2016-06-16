package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
)

type HoverflyDirectory struct {
	Path string
}

func NewHoverflyDirectory() (HoverflyDirectory) {
	return HoverflyDirectory{Path: "/"}
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