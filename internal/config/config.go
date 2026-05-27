package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetDir          string   `yaml:"target_dir"`
	Extensions         []string `yaml:"extensions"`
	RetentionDays      int      `yaml:"retention_days"`
	DryRun             bool     `yaml:"dry_run"`
	WarningAck         bool     `yaml:"warning_ack"`
	LogPath            string   `yaml:"log_path"`
	CheckIntervalHours int      `yaml:"check_interval_hours"`
}

func Load() (*Config, error) {
	path := ConfigPath()
	// #nosec G304 — path is an internal config path, not user input.
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func Reset() error {
	path := ConfigPath()
	return os.Remove(path)
}

func Exists() bool {
	_, err := os.Stat(ConfigPath())
	return err == nil
}
