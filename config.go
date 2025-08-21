package main

import (
	"fmt"
	"github.com/cloudogu/github-forgejo-backup/internal/logs"
	"github.com/goccy/go-yaml"
	"os"
	"path/filepath"
)

type Config struct {
	GithubBaseUrl  string `yaml:"github_base_url"`
	GithubOrga     string `yaml:"github_orga"`
	ForgejoBaseUrl string `yaml:"forgejo_base_url"`
	ForgejoOrga    string `yaml:"forgejo_orga"`
	WebhoolUrl     string `yaml:"webhook_url"`
	CronSpec       string `yaml:"cron-spec"`
	TimeZone       string `yaml:"time-zone"`
}

type Tokens struct {
	Github  string `yaml:"github"`
	Forgejo string `yaml:"forgejo"`
}

var config Config
var tokens Tokens

func (c Config) Load() error {

	if err := loadConfig(); err != nil {
		return err
	}

	if err := loadTokens(); err != nil {
		return err
	}

	return nil
}

func loadConfig() error {

	path, err := filepath.Abs("config.yaml")
	if err != nil {
		return fmt.Errorf("failed resolving config file path: %v\n", err.Error())
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

func loadTokens() error {

	path, err := filepath.Abs("tokens.yaml")
	if err != nil {
		return fmt.Errorf("failed resolving tokens file path: %v\n", err.Error())
	}

	logs.Info("loading tokens", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading tokens file: %v\n", err.Error())
	}

	err = yaml.Unmarshal(data, &tokens)
	if err != nil {
		return fmt.Errorf("failed unmarshalling tokens data: %v\n", err.Error())
	}

	return nil
}
