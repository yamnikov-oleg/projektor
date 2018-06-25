package main

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type ProjektorConfig struct {
	KeyBind            string
	ForceCacheInterval int64
	EnabledCategories  struct {
		Calc      bool
		History   bool
		Apps      bool
		URL       bool
		Commands  bool
		Files     bool
		WebSearch bool
	}
	History struct {
		Capacity int
	}
	URL struct {
		Icon string
	}
	WebSearch struct {
		Engine string
		Icon   string
	}
	UI struct {
		Width  string
		Height string
	}
}

var (
	ConfigFilePath                  = path.Join(AppDir, "config.yaml")
	Config         *ProjektorConfig = MustLoadConfig()
)

func DefaultConfig() *ProjektorConfig {
	c := &ProjektorConfig{}

	c.KeyBind = "mod4-q"
	c.ForceCacheInterval = 60

	c.EnabledCategories.Calc = true
	c.EnabledCategories.History = true
	c.EnabledCategories.Apps = true
	c.EnabledCategories.URL = true
	c.EnabledCategories.Commands = true
	c.EnabledCategories.Files = true
	c.EnabledCategories.WebSearch = true

	c.History.Capacity = 40

	c.URL.Icon = "web-browser"

	c.WebSearch.Engine = "https://google.com/search?q=%s"
	c.WebSearch.Icon = "web-browser"

	c.UI.Width = "400"
	c.UI.Height = "480"

	return c
}

func SaveConfig(cfg *ProjektorConfig) error {
	err := os.MkdirAll(AppDir, 0700)
	if err != nil {
		return err
	}

	f, err := os.Create(ConfigFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
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

	config := DefaultConfig()
	err = yaml.Unmarshal(contents, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfig() *ProjektorConfig {
	var config *ProjektorConfig
	var err error

	config, err = OpenConfig()
	if err != nil {
		errduring("opening config file at %q", err, "Attempting to create one", ConfigFilePath)
		config = DefaultConfig()
	}

	if err := SaveConfig(config); err != nil {
		errduring("creating config file at %q", err, "Using default options", ConfigFilePath)
		return config
	}

	return config
}
