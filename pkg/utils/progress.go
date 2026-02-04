package utils

import (
	"fmt"
	"io"
)

// ProgressReader wraps an io.Reader to track progress
type ProgressReader struct {
	reader   io.Reader
	total    int64
	current  int64
	callback func(current, total int64)
}

// NewProgressReader creates a new progress reader
func NewProgressReader(reader io.Reader, total int64, callback func(current, total int64)) *ProgressReader {
	return &ProgressReader{
		reader:   reader,
		total:    total,
		callback: callback,
	}
}

// Read implements io.Reader
func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.current += int64(n)
	if pr.callback != nil {
		pr.callback(pr.current, pr.total)
	}
	return n, err
}

// ProgressBar represents a simple progress bar
type ProgressBar struct {
	total   int64
	current int64
	width   int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		total: total,
		width: 50,
	}
}

// Update updates the progress
func (pb *ProgressBar) Update(current int64) {
	pb.current = current
	percent := float64(current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * float64(current) / float64(pb.total))
	bar := "["
	for i := 0; i < pb.width; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"
	fmt.Printf("\r%s %.2f%%", bar, percent)
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	fmt.Println()
}
