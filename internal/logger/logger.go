package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const LogTimeFormat = "2006/01/02 15:04:05"

// Prune removes log entries older than retentionDays from path.
// Lines that do not start with a parseable timestamp are preserved.
func Prune(path string, retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if info.Size() == 0 {
		return nil
	}

	cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)

	// #nosec G304 — path is an internal log path, not user input.
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	tmp := path + ".tmp"
	// #nosec G304 — tmp is derived from the internal log path, not user input.
	dst, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dst.Close()

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		ts, ok := parseLogTimestamp(line)
		if ok && ts.Before(cutoff) {
			continue
		}
		if _, err := fmt.Fprintln(dst, line); err != nil {
			// #nosec G104 — best-effort cleanup of temporary file.
			os.Remove(tmp)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		// #nosec G104 — best-effort cleanup of temporary file.
		os.Remove(tmp)
		return err
	}

	if err := dst.Close(); err != nil {
		// #nosec G104 — best-effort cleanup of temporary file.
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, path)
}

func parseLogTimestamp(line string) (time.Time, bool) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return time.Time{}, false
	}
	ts, err := time.Parse(LogTimeFormat, parts[0]+" "+parts[1])
	if err != nil {
		return time.Time{}, false
	}
	return ts, true
}
