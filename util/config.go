package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config is the config.json but as a struct
var Config *configuration

type configuration struct {
	Token  string `json:"token"` 
	Prefix string `json:"prefix"`
}


func init() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Trouble opening the config: %s", err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Trouble parsing the config: %s", err)
	}

	return
}
