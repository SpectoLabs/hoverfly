package main

import (
	"strings"
	"errors"
	"fmt"
)

type Hoverfile struct {
	Vendor  string
	Name    string
	Version string
}

func (h *Hoverfile) GetFileName() string {
	return fmt.Sprintf("%v.%v.%v.hfile", h.Vendor, h.Name, h.Version)
}

func (h *Hoverfile) String() string {
	return h.Vendor + "/" + h.Name + ":" + h.Version
}

func NewHoverfile(key string) (Hoverfile, error) {

	var vendor, name, version string

	if strings.ContainsAny(key, "@\\~`|[]{}!@€£$%^&*()+=?\"';§±") {
		return Hoverfile{}, errors.New("Invalid characters used in hoverfile name")
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
		version = "v1"
	} else if !strings.Contains(key, "/") && ! strings.Contains(key, ":") {
		vendor = ""
		name = key
		version = "v1"
	}

	return Hoverfile {
		Vendor: vendor,
		Name: name,
		Version: version,
	}, nil
}