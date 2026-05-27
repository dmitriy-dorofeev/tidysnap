package config

import (
	"errors"
	"os"
	"path/filepath"
)

func homeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

func DefaultConfig() (*Config, error) {
	home, err := homeDir()
	if err != nil {
		return nil, err
	}
	return &Config{
		TargetDir:          home,
		Extensions:         []string{".png", ".jpg", ".jpeg", ".mov", ".mp4", ".gif"},
		RetentionDays:      30,
		DryRun:             true,
		WarningAck:         false,
		LogPath:            filepath.Join(home, "Library", "Logs", "tidysnap", "cleanup.log"),
		CheckIntervalHours: 24,
	}, nil
}

func ConfigDir() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Application Support", "tidysnap"), nil
}

func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func LogDir() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Logs", "tidysnap"), nil
}

func LogPath() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "Logs", "tidysnap", "cleanup.log"), nil
}

func PlistPath() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library", "LaunchAgents", "com.tidysnap.plist"), nil
}

func Cleanup() error {
	var errs []error

	if dir, err := ConfigDir(); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, err)
		}
	}

	if dir, err := LogDir(); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, err)
		}
	}

	if path, err := PlistPath(); err == nil {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
