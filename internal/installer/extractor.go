package installer

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor handles archive extraction
type Extractor struct{}

// NewExtractor creates a new Extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts an archive to the destination directory
func (e *Extractor) Extract(archivePath, destPath string) error {
	ext := strings.ToLower(filepath.Ext(archivePath))
	
	switch ext {
	case ".zip":
		return e.extractZip(archivePath, destPath)
	case ".tar", ".gz", ".tgz":
		return e.extractTar(archivePath, destPath)
	case ".exe":
		// For .exe installers, we might just move them or run them
		// For now, just copy to destination
		return e.copyFile(archivePath, filepath.Join(destPath, filepath.Base(archivePath)))
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// extractZip extracts a ZIP archive
func (e *Extractor) extractZip(archivePath, destPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if err := e.extractZipFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile extracts a single file from a ZIP archive
func (e *Extractor) extractZipFile(f *zip.File, destPath string) error {
	// Clean the file path to prevent zip slip
	fpath := filepath.Join(destPath, f.Name)
	if !strings.HasPrefix(fpath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", f.Name)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(fpath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

// extractTar extracts a TAR archive (including .tar.gz, .tgz)
func (e *Extractor) extractTar(archivePath, destPath string) error {
	// This is a placeholder - in production, you'd use archive/tar
	// and compress/gzip packages
	return fmt.Errorf("TAR extraction not yet implemented")
}

// copyFile copies a file from src to dst
func (e *Extractor) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
