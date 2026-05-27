package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

const label = "com.tidysnap"

func Install(binaryPath string, intervalHours int) error {
	plist := GeneratePlist(label, binaryPath, intervalHours)
	if err := WritePlist(plist); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}

	plistPath, err := config.PlistPath()
	if err != nil {
		return fmt.Errorf("plist path: %w", err)
	}
	// #nosec G204 — plistPath is an internal system path, not user input.
	if err := exec.Command("launchctl", "load", plistPath).Run(); err != nil {
		return fmt.Errorf("launchctl load: %w", err)
	}
	if err := exec.Command("launchctl", "start", label).Run(); err != nil {
		return fmt.Errorf("launchctl start: %w", err)
	}
	return nil
}

func Uninstall() error {
	var errs []error

	_ = exec.Command("launchctl", "stop", label).Run()

	plistPath, err := config.PlistPath()
	if err == nil {
		// #nosec G204 — plistPath is an internal system path, not user input.
		_ = exec.Command("launchctl", "unload", plistPath).Run()
	}

	if err := RemovePlist(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func IsInstalled() bool {
	plistPath, err := config.PlistPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(plistPath)
	return err == nil
}

func IsLoaded() bool {
	_, err := exec.Command("launchctl", "list", label).Output()
	return err == nil
}

func IsRunning() bool {
	out, err := exec.Command("launchctl", "list").Output()
	if err != nil {
		return false
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[2] == label {
			return fields[0] != "-" && fields[0] != ""
		}
	}
	return false
}

func Load() error {
	plistPath, err := config.PlistPath()
	if err != nil {
		return err
	}
	// #nosec G204 — plistPath is an internal system path, not user input.
	return exec.Command("launchctl", "load", plistPath).Run()
}

func Unload() error {
	plistPath, err := config.PlistPath()
	if err != nil {
		return err
	}
	// #nosec G204 — plistPath is an internal system path, not user input.
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

func NextRunTime(intervalHours int) (time.Time, bool) {
	if !IsInstalled() || !IsLoaded() {
		return time.Time{}, false
	}
	logPath, err := config.LogPath()
	if err != nil {
		return time.Time{}, false
	}
	info, err := os.Stat(logPath)
	if err != nil {
		return time.Time{}, false
	}
	nextRun := info.ModTime().Add(time.Duration(intervalHours) * time.Hour)
	return nextRun, true
}
