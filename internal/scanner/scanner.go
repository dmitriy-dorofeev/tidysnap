package scanner

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type ScanResult struct {
	Path    string
	Size    int64
	ModTime time.Time
	Ext     string
}

type CleanupStats struct {
	FilesRemoved int
	BytesFreed   int64
	Errors       []string
	Timestamp    time.Time
}

type Scanner struct {
	extensions []string
	retention  time.Duration
}

func New(extensions []string, retentionDays int) *Scanner {
	lowers := make([]string, len(extensions))
	for i, e := range extensions {
		lowers[i] = strings.ToLower(e)
	}
	return &Scanner{
		extensions: lowers,
		retention:  time.Duration(retentionDays) * 24 * time.Hour,
	}
}

func (s *Scanner) Scan(targetDir string) ([]ScanResult, error) {
	var results []ScanResult

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if path == targetDir {
				return err
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !slices.Contains(s.extensions, ext) {
			return nil
		}

		if time.Since(info.ModTime()) <= s.retention {
			return nil
		}

		results = append(results, ScanResult{
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Ext:     ext,
		})

		return nil
	})

	return results, err
}
