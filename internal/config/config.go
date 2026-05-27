package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Language           string   `yaml:"language"`
	TargetDir          string   `yaml:"target_dir"`
	Extensions         []string `yaml:"extensions"`
	RetentionDays      int      `yaml:"retention_days"`
	DryRun             bool     `yaml:"dry_run"`
	WarningAck         bool     `yaml:"warning_ack"`
	LogPath            string   `yaml:"log_path"`
	CheckIntervalHours int      `yaml:"check_interval_hours"`
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	// #nosec G304 — path is an internal config path, not user input.
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig()
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
	path, err := ConfigPath()
	if err != nil {
		return err
	}
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
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func Exists() bool {
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
