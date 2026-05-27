package cleaner

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
	"github.com/dustin/go-humanize"
)

type Cleaner struct {
	dryRun bool
	logger *log.Logger
}

func New(dryRun bool, logger *log.Logger) *Cleaner {
	return &Cleaner{dryRun: dryRun, logger: logger}
}

func (c *Cleaner) Clean(files []scanner.ScanResult) (*scanner.CleanupStats, error) {
	stats := &scanner.CleanupStats{Timestamp: time.Now()}

	for _, file := range files {
		if c.dryRun {
			c.logger.Printf("[DRY RUN] Would delete: %s (%s)", file.Path, humanize.Bytes(uint64(file.Size)))
			stats.FilesRemoved++
			stats.BytesFreed += file.Size
			continue
		}

		if err := os.Remove(file.Path); err != nil {
			stats.Errors = append(stats.Errors, fmt.Sprintf("%s: %v", file.Path, err))
			continue
		}

		stats.FilesRemoved++
		stats.BytesFreed += file.Size
		c.logger.Printf("Deleted: %s (%s)", file.Path, humanize.Bytes(uint64(file.Size)))
	}

	return stats, nil
}
