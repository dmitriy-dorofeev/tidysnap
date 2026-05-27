package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain_Version(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Build the binary first
		return
	}
	cmd := exec.Command("go", "run", ".", "-version")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version flag failed: %v, output: %s", err, out)
	}
	if !strings.Contains(string(out), "tidysnap") {
		t.Errorf("output should contain tidysnap, got: %s", out)
	}
}

func TestMain_ConfigPath(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-config-path")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config-path flag failed: %v, output: %s", err, out)
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		t.Error("config-path should not be empty")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("basename = %q, want config.yaml", filepath.Base(path))
	}
}

func TestMain_Reset(t *testing.T) {
	tmp, err := os.MkdirTemp("", "tidysnap-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	// Create a config file
	cfgPath := filepath.Join(tmp, "Library", "Application Support", "tidysnap", "config.yaml")
	os.MkdirAll(filepath.Dir(cfgPath), 0755)
	os.WriteFile(cfgPath, []byte("target_dir: /tmp\n"), 0644)

	cmd := exec.Command("go", "run", ".", "-reset")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("reset flag failed: %v, output: %s", err, out)
	}
	if !strings.Contains(string(out), "сброшены") {
		t.Errorf("output should indicate reset, got: %s", out)
	}

	if _, err := os.Stat(cfgPath); !os.IsNotExist(err) {
		t.Error("config file should be removed after reset")
	}
}

func configPathForTest(t *testing.T) string {
	t.Helper()
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "Application Support", "tidysnap", "config.yaml")
}
