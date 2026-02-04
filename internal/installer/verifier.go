package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Verifier handles checksum verification
type Verifier struct{}

// NewVerifier creates a new verifier
func NewVerifier() *Verifier {
	return &Verifier{}
}

// VerifyChecksum verifies the checksum of a file
func (v *Verifier) VerifyChecksum(filePath, expectedChecksum string) error {
	if expectedChecksum == "" {
		return nil // Skip verification if no checksum provided
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}
