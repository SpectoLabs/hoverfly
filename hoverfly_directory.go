package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"os"
	"path"
	"fmt"
	"io/ioutil"
	"strconv"
	"errors"
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

func (h *HoverflyDirectory) GetPid(adminPort, proxyPort string) (int) {
	hoverflyPidFile := h.buildPidFilePath(adminPort, proxyPort)
	if fileIsPresent(hoverflyPidFile) {
		pidFileData, err := ioutil.ReadFile(hoverflyPidFile)
		if err != nil {
			log.Debug(err.Error())
			log.Fatal("Could not get pid from .hoverfly directory")
		}

		pid, err := strconv.Atoi(string(pidFileData))
		if err != nil {
			log.Debug(err.Error())
			log.Fatal("Could not read pid from .hoverfly directory")
		}

		return pid
	}

	return 0
}

func (h *HoverflyDirectory) WritePid(adminPort, proxyPort string, pid int) (error) {
	pidFilePath := h.buildPidFilePath(adminPort, proxyPort)
	if fileIsPresent(pidFilePath) {
		return errors.New("Hoverfly pid already exists")
	}
	return ioutil.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0644)
}

func (h *HoverflyDirectory) DeletePid(adminPort, proxyPort string) (error) {
	return os.Remove(h.buildPidFilePath(adminPort, proxyPort))
}

func (h *HoverflyDirectory) buildPidFilePath(adminPort, proxyPort string) (string) {
	pidName := fmt.Sprintf("hoverfly.%v.%v.pid", adminPort, proxyPort)
	return filepath.Join(h.Path, pidName)
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