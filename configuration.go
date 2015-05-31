package main

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	ApiToken	string
}

func (config *Configuration) Load() bool {
	f, err := ioutil.ReadFile(config.filename())

	if err != nil {
		return false
	}

	yaml.Unmarshal(f, &config)

	return true
}

func (config *Configuration) Save() {
	f, _ := os.Create(config.filename())
	defer f.Close()

	b, _ := yaml.Marshal(config)
	f.Write(b)
}

func (config *Configuration) filename() string {
	return os.Getenv("HOME") + "/.togglrc"
}