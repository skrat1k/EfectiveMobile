package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env-required:"true"`
	Database `yaml:"database"`
}

type Database struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

func MustLoad() (*Config, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get work directory: %w", err)
	}

	configPath, err := filepath.Abs(filepath.Join(workdir, "..", "config", "local.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file dots not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	return &cfg, nil
}
