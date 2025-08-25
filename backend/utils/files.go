package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// UploadDir returns the directory used to store uploaded files.
// Controlled by UPLOAD_DIR env var; must be set. The process exits if missing.
func UploadDir() string {
	if d := os.Getenv("UPLOAD_DIR"); d != "" {
		return d
	}
	log.Fatal("UPLOAD_DIR environment variable must be set")
	return "" // unreachable
}

// UploadImage saves an uploaded image to disk and returns the public URL path (e.g., /files/name.ext).
// Parameters:
// - baseName: logical name to use for the filename base (will be sanitized)
// - originalFilename: the original uploaded filename (used to derive extension if present)
// - contentType: detected MIME type (used to derive extension if original missing)
// - head: the first bytes already read from the file stream (used for content that was sniffed)
// - n: number of valid bytes in head
// - rest: an io.Reader for the remaining file content to write
func UploadImage(baseName, originalFilename, contentType string, head []byte, n int, rest io.Reader) (string, error) {
	// Ensure images subdirectory under the configured upload dir exists (do not create root here)
	imagesDir := filepath.Join(UploadDir(), "images")
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		return "", fmt.Errorf("prepare images dir: %w", err)
	}
	// determine extension
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		case "image/gif":
			ext = ".gif"
		default:
			return "", fmt.Errorf("unknown file type")
		}
	}

	base := sanitizeFileBase(baseName)
	if base == "" {
		base = "file"
	}
	fname := base + ext
	diskPath := filepath.Join(imagesDir, fname)
	// avoid collisions
	if _, err := os.Stat(diskPath); err == nil {
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
			candidatePath := filepath.Join(imagesDir, candidate)
			if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
				fname = candidate
				diskPath = candidatePath
				break
			}
		}
	}

	dst, err := os.Create(diskPath)
	if err != nil {
		return "", fmt.Errorf("save file: %w", err)
	}
	defer dst.Close()

	if n > 0 {
		if _, err := dst.Write(head[:n]); err != nil {
			return "", fmt.Errorf("write head: %w", err)
		}
	}
	if _, err := io.Copy(dst, rest); err != nil {
		return "", err
	}
	return "/files/images/" + fname, nil
}

// sanitizeFileBase converts a name into a safe filename base.
func sanitizeFileBase(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return out
	}
	var b2 strings.Builder
	prevDash := false
	for _, r := range out {
		if r == '-' {
			if !prevDash {
				b2.WriteRune(r)
			}
			prevDash = true
			continue
		}
		prevDash = false
		b2.WriteRune(r)
	}
	return b2.String()
}

// VerifyUploadDirWritable checks that UPLOAD_DIR exists, is a directory, and is writable.
// Returns an error if any condition is not met. This does not create the directory.
func VerifyUploadDirWritable() error {
	dir := UploadDir() // exits if env missing
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("upload dir %q does not exist", dir)
		}
		return fmt.Errorf("cannot access upload dir %q: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("upload dir %q is not a directory", dir)
	}
	// probe writability by creating a temp file
	f, err := os.CreateTemp(dir, ".probe-*")
	if err != nil {
		return fmt.Errorf("upload dir %q is not writable: %w", dir, err)
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
	return nil
}
