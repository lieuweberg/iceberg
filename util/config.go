package util

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the config.json but as a struct
var Config *configuration

type configuration struct {
	Token string `json:"token"`
}

func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		panic(err)
	}

	return
}
