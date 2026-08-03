package fsutil

import (
	"fmt"
	"os"
)

// Exists returns true if the file or directory at path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WriteNew writes content to path only if the file does not already exist.
func WriteNew(path string, content []byte) error {
	if Exists(path) {
		return fmt.Errorf("file already exists: %s", path)
	}
	return os.WriteFile(path, content, 0o644)
}

// WriteWithBackup writes content to path. If path exists, it copies the existing
// file to path + ".bak" before overwriting it.
func WriteWithBackup(path string, content []byte) error {
	if Exists(path) {
		old, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path+".bak", old, 0o644); err != nil {
			return err
		}
	}
	return os.WriteFile(path, content, 0o644)
}
