package config

import (
	"os"
	"path/filepath"
)

func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		TargetDir:          home,
		Extensions:         []string{".png", ".jpg", ".jpeg", ".mov", ".mp4", ".gif"},
		RetentionDays:      30,
		DryRun:             true,
		WarningAck:         false,
		LogPath:            filepath.Join(home, "Library", "Logs", "tidysnap", "cleanup.log"),
		CheckIntervalHours: 24,
	}
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "tidysnap")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func LogDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Logs", "tidysnap")
}

func LogPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Logs", "tidysnap", "cleanup.log")
}

func PlistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", "com.tidysnap.plist")
}
