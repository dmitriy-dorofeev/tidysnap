package cleaner

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
	"github.com/dustin/go-humanize"
)

type CleanupStats struct {
	FilesRemoved int
	BytesFreed   int64
	Errors       []string
	Timestamp    time.Time
}

type Cleaner struct {
	dryRun bool
	logger *log.Logger
}

func New(dryRun bool, logger *log.Logger) *Cleaner {
	return &Cleaner{dryRun: dryRun, logger: logger}
}

func safeUint64(n int64) uint64 {
	if n < 0 {
		return 0
	}
	return uint64(n)
}

func (c *Cleaner) Clean(ctx context.Context, files []scanner.ScanResult) (*CleanupStats, error) {
	stats := &CleanupStats{Timestamp: time.Now()}

	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return stats, err
		}

		if c.dryRun {
			c.logger.Printf("[DRY RUN] Would delete: %s (%s)", file.Path, humanize.Bytes(safeUint64(file.Size)))
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
		c.logger.Printf("Deleted: %s (%s)", file.Path, humanize.Bytes(safeUint64(file.Size)))
	}

	return stats, nil
}
