package main

import (
	"io/ioutil"

	"github.com/thrashy/drum-machine-808/model"
)

// LoadConfigurationFromToml takes a file name containing toml configuration and returns a structured response
// containing configuration information for the application.
func loadConfigurationFromToml(fileName string) (*model.SongBeatConfiguration, error) {

	tomlData, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	sb := model.SongBeatConfiguration{}
	err = sb.UnmarshalTOML(tomlData)

	return &sb, err
}
