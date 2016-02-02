package config

import (
	"log"
    "io/ioutil"

	"github.com/hashicorp/hcl"
)

type Config struct {
	IgnoreText []string
	IgnoreSource []string
}

func GetConfig(configFile string) Config {
	var conf Config

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error: Couldn't read config file %s: %s", configFile, err)
	}

	err = hcl.Decode(&conf,string(configData))
	if err != nil {
		log.Fatalf("Error parsing config file %s: %s", configFile, err)
	}

	return conf
}
