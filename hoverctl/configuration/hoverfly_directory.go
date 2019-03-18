package configuration

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
)

type HoverflyDirectory struct {
	Path string
}

func NewHoverflyDirectory(config Config) (HoverflyDirectory, error) {
	var hoverflyDirectory HoverflyDirectory

	if len(config.GetFilepath()) == 0 {
		log.Debug("Missing a config file")
		log.Debug("Creating a new  a config file")
		hoverflyDirectory = HoverflyDirectory{Path: createHoverflyDirectory(getHomeDirectory())}

		err := config.WriteToFile(hoverflyDirectory)
		if err != nil {
			log.Debug(err.Error())
			return HoverflyDirectory{}, errors.New("Could not write new config to disk")
		}

	} else {
		hoverflyDirectory = HoverflyDirectory{
			Path: filepath.Dir(config.GetFilepath()),
		}
	}

	return hoverflyDirectory, nil
}

func getHomeDirectory() string {
	homeDirectory, err := homedir.Dir()
	if err != nil {
		log.Debug(err.Error())
		log.Fatal("Unable to get home directory")
	}

	return homeDirectory
}

func createHoverflyDirectory(homeDirectory string) string {
	hoverflyDirectory := filepath.Join(homeDirectory, "/.hoverfly")

	if !fileIsPresent(hoverflyDirectory) {
		err := os.Mkdir(hoverflyDirectory, 0777)
		if err != nil {
			log.Debug(err.Error())
			log.Fatal("Could not create a .hoverfly directory")
		}

		return hoverflyDirectory
	}

	return hoverflyDirectory
}

func fileIsPresent(fileURI string) bool {
	if _, err := os.Stat(fileURI); err != nil {
		return os.IsExist(err)
	}

	return true
}
