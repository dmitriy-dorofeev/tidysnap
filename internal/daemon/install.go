package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

const label = "com.tidysnap"

func Install(binaryPath string, intervalHours int) error {
	plist := GeneratePlist(label, binaryPath, intervalHours)
	if err := WritePlist(plist); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}

	plistPath := config.PlistPath()
	if err := exec.Command("launchctl", "load", plistPath).Run(); err != nil {
		return fmt.Errorf("launchctl load: %w", err)
	}
	if err := exec.Command("launchctl", "start", label).Run(); err != nil {
		return fmt.Errorf("launchctl start: %w", err)
	}
	return nil
}

func Uninstall() error {
	plistPath := config.PlistPath()

	_ = exec.Command("launchctl", "stop", label).Run()
	_ = exec.Command("launchctl", "unload", plistPath).Run()

	_ = RemovePlist()
	return nil
}

func IsInstalled() bool {
	plistPath := config.PlistPath()
	_, err := os.Stat(plistPath)
	return err == nil
}

func BinaryPath() string {
	ex, err := os.Executable()
	if err != nil {
		return filepath.Join("/usr/local/bin", "tidysnap")
	}
	return ex
}
