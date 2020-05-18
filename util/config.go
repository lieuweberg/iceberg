package util

import (
	"encoding/json"
	"io/ioutil"
)

// Configuration The config.json file
type Configuration struct {
	Token string `json:"token"`
}

// Config Returns a Configuration struct and an error if there was one
func Config() (configuration Configuration, err error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &configuration)

	return
}
