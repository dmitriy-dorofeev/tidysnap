package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func IsLoaded() bool {
	out, err := exec.Command("launchctl", "list", label).Output()
	if err != nil {
		return false
	}
	fields := strings.Fields(string(out))
	return len(fields) >= 3 && fields[2] == label
}

func IsRunning() bool {
	out, err := exec.Command("launchctl", "list", label).Output()
	if err != nil {
		return false
	}
	fields := strings.Fields(string(out))
	if len(fields) >= 3 && fields[2] == label {
		return fields[0] != "-" && fields[0] != ""
	}
	return false
}

func Load() error {
	plistPath := config.PlistPath()
	return exec.Command("launchctl", "load", plistPath).Run()
}

func Unload() error {
	plistPath := config.PlistPath()
	return exec.Command("launchctl", "unload", plistPath).Run()
}

func Stop() error {
	return exec.Command("launchctl", "stop", label).Run()
}

func Start() error {
	return exec.Command("launchctl", "start", label).Run()
}

func BinaryPath() string {
	ex, err := os.Executable()
	if err != nil {
		return filepath.Join("/usr/local/bin", "tidysnap")
	}
	return ex
}
