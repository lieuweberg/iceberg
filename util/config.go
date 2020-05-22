package util

import (
	"encoding/json"
	"io/ioutil"
)

// Configuration is the config.json file, but as a struct
type Configuration struct {
	Token string `json:"token"`
}

// LoadConfig returns a Configuration struct and an error if there was one
func LoadConfig() (configuration Configuration, err error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &configuration)

	return
}
