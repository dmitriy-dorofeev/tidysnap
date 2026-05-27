package scanner

import (
	"context"
	"io/fs"
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

func (s *Scanner) Scan(ctx context.Context, targetDir string) ([]ScanResult, error) {
	var results []ScanResult

	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == targetDir {
				return err
			}
			return nil
		}

		if err := ctx.Err(); err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
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
