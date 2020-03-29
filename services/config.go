package services

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramToken string `yaml:"telegram_token"`
}

func GetConfig() (Config, error) {
	f, err := os.Open("config.yml")

	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
