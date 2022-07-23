package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Config is the Mokou configuration struct
type Config struct {
	InitialNap    string
	EienteiConfig *EienteiConfig
	AsagiConfig   *AsagiConfig
	BatchSize     int
	Boards        []BoardConfig
}

//EienteiConfig is the configuration struct for
//connection to an Eientei db instance.
type EienteiConfig struct {
	ConnectionString string
}

//AsagiConfig is the configuration struct for
//connection to an AsagiConfig db instance.
type AsagiConfig struct {
	ConnectionString string
}

//BoardConfig is the configuration struct for
//some board to be imported.
type BoardConfig struct {
	Name          string
	EnableCode    bool
	EnableSpoiler bool
	EnableFortune bool
	EnableExif    bool
	EnableOekaki  bool
}

//LoadConfig reads config.json and unmarshals it into a Config struct.
//Errors might be returned due to IO or invalid JSON.
func LoadConfig() (Config, error) {
	blob, err := ioutil.ReadFile("config.json")

	if err != nil {
		return Config{}, fmt.Errorf("Error loading file config.json in project root: %s", err)
	}

	var conf Config

	err = json.Unmarshal(blob, &conf)

	if err != nil {
		return Config{}, fmt.Errorf(
			"Error unmarshalling configuration file contents to JSON:\n File contents: %s\n Error message: %s",
			blob,
			err,
		)
	}

	return conf, nil
}
