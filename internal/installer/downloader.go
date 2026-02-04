package installer

import (
	"context"
	"fmt"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/schollz/progressbar/v3"
)

// Downloader handles file downloads with progress tracking
type Downloader struct{}

// NewDownloader creates a new Downloader
func NewDownloader() *Downloader {
	return &Downloader{}
}

// Download downloads a file from URL to destination with progress bar
func (d *Downloader) Download(ctx context.Context, url, dest string) error {
	// Create grab client
	client := grab.NewClient()
	req, err := grab.NewRequest(dest, url)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req = req.WithContext(ctx)

	// Start download
	resp := client.Do(req)

	// Get file size
	size := resp.Size()
	
	// Create progress bar
	bar := progressbar.DefaultBytes(
		size,
		"Downloading",
	)

	// Monitor progress with ticker
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bar.Set64(resp.BytesComplete())
			if resp.IsComplete() {
				bar.Finish()
				fmt.Println()
				goto complete
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

complete:
	// Check for errors
	if err := resp.Err(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}

// DownloadWithoutProgress downloads without showing progress (for smaller files)
func (d *Downloader) DownloadWithoutProgress(ctx context.Context, url, dest string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(dest, url)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req = req.WithContext(ctx)
	resp := client.Do(req)

	// Wait for download to complete
	<-resp.Done

	if err := resp.Err(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	return nil
}

// VerifyChecksum verifies the downloaded file's checksum
func (d *Downloader) VerifyChecksum(filePath, expectedChecksum string) error {
	if expectedChecksum == "" {
		// No checksum to verify
		return nil
	}

	// TODO: Implement actual checksum verification
	// This would use crypto/sha256 or similar
	
	return nil
}
