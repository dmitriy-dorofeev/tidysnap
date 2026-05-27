package daemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func GeneratePlist(label, binaryPath string, intervalHours int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
        <string>--cleanup</string>
    </array>
    <key>StartInterval</key>
    <integer>%d</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>%s</string>
    <key>StandardErrorPath</key>
    <string>%s</string>
</dict>
</plist>`,
		label,
		binaryPath,
		intervalHours*3600,
		filepath.Join(os.Getenv("HOME"), "Library", "Logs", "tidysnap", "stdout.log"),
		filepath.Join(os.Getenv("HOME"), "Library", "Logs", "tidysnap", "stderr.log"),
	)
}

func WritePlist(content string) error {
	path := config.PlistPath()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0600)
}

func RemovePlist() error {
	return os.Remove(config.PlistPath())
}
