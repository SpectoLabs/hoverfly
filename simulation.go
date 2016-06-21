package main

import (
	"strings"
	"errors"
	"fmt"
)

type Simulation struct {
	Vendor  string
	Name    string
	Version string
}

func (s *Simulation) GetFileName() string {
	return fmt.Sprintf("%v.%v.%v.json", s.Vendor, s.Name, s.Version)
}

func (s *Simulation) String() string {
	return s.Vendor + "/" + s.Name + ":" + s.Version
}

func NewSimulation(key string) (Simulation, error) {

	var vendor, name, version string

	if strings.ContainsAny(key, "@\\~`|[]{}!@€£$%^&*()+=?\"';§±") {
		return Simulation{}, errors.New("Invalid characters used in simulation name")
	}

	if strings.Contains(key, "/") && strings.Contains(key, ":") {
		vendorAndEnd := strings.Split(key, "/", )
		nameAndVersion := strings.Split(vendorAndEnd[1], ":")

		vendor = vendorAndEnd[0]
		name = nameAndVersion[0]
		version = nameAndVersion[1]
	} else if strings.Contains(key,"/") {
		vendorAndEnd := strings.Split(key, "/", )
		vendor = vendorAndEnd[0]
		name = vendorAndEnd[1]
		version = "latest"
	} else if !strings.Contains(key, "/") && ! strings.Contains(key, ":") {
		vendor = ""
		name = key
		version = "latest"
	}

	return Simulation{
		Vendor: vendor,
		Name: name,
		Version: version,
	}, nil
}