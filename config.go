package main

import (
	"fmt"
	"os"
	"path/filepath"
)

import "github.com/goccy/go-yaml"

type Config struct {
	GithubBaseUrl  string `yaml:"github_base_url"`
	GithubOrga     string `yaml:"github_orga"`
	GithubToken    string `yaml:"github_token"`
	ForgejoBaseUrl string `yaml:"forgejo_base_url"`
	ForgejoOrga    string `yaml:"forgejo_orga"`
	ForgejoToken   string `yaml:"forgejo_token"`
	CronSpec       string `yaml:"cron-spec"`
	TimeZone       string `yaml:"time-zone"`
}

var config Config

func (c Config) Load() error {

	path, err := filepath.Abs("config.yaml")
	if err != nil {
		return fmt.Errorf("failed resolving config path: %v\n", err.Error())
	}

	logs.Info("loading config", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading config file: %v\n", err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("failed unmarshalling config data: %v\n", err.Error())
	}

	return nil

}
