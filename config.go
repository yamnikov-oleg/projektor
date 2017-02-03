package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type ProjektorConfig struct {
	EnabledCategories struct {
		History  bool
		Apps     bool
		URL      bool
		Commands bool
		Files    bool
	}
	History struct {
		Capacity int
	}
}

var (
	ConfigFilePath                  = path.Join(AppDir, "config.yaml")
	Config         *ProjektorConfig = MustLoadConfig()
)

func DefaultConfig() *ProjektorConfig {
	c := &ProjektorConfig{}

	c.EnabledCategories.History = true
	c.EnabledCategories.Apps = true
	c.EnabledCategories.URL = true
	c.EnabledCategories.Commands = true
	c.EnabledCategories.Files = true

	c.History.Capacity = 40

	return c
}

func OpenConfig() (*ProjektorConfig, error) {
	f, err := os.Open(ConfigFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	config := &ProjektorConfig{}
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func CreateConfig() (*ProjektorConfig, error) {
	err := os.MkdirAll(AppDir, 0700)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(ConfigFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := DefaultConfig()
	buf, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(buf)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfig() *ProjektorConfig {
	config, err := OpenConfig()
	if err == nil {
		return config
	}

	errduring("opening config file at %q", err, "Attempting to create one", ConfigFilePath)
	config, err = CreateConfig()
	if err == nil {
		return config
	}

	errduring("creating config file at %q", err, "Using default options", ConfigFilePath)
	return DefaultConfig()
}
