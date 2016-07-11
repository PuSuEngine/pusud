package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"fmt"
)

var settingsFilename = "settings.yaml"
var settingsContents *[]byte = nil
var DEBUG = false

type Settings struct {
	Authenticator   string   `yaml:"authenticator"`
	Relays          []string `yaml:"relays"`
	ClientPort      int      `yaml:"client_port"`
	NetworkPort     int      `yaml:"network_port"`
	AllowedChannels []string `yaml:"allowed_channels"`
}

func GetSettingsContents() *[]byte {
	if settingsContents == nil {
		contents, err := ioutil.ReadFile(settingsFilename)

		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("Settings file %s not found. Try copying settings.example.yaml", settingsFilename)
			} else {
				log.Fatalf("Error reading settings: %v", err)
			}
		}

		log.Printf("Read settings from %s", settingsFilename)

		settingsContents = &contents
	}

	return settingsContents
}

func GetSettings() *Settings {
	s := Settings{}
	data := GetSettingsContents()
	yaml.Unmarshal(*data, &s)

	if DEBUG {
		fmt.Printf("%+v\n", s)
	}

	return &s
}
