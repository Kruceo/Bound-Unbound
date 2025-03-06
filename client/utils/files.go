package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func CombinedFileHash(files []string) (string, error) {
	hash := sha256.New()
	for _, filePath := range files {

		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("error to open file %s: %w", filePath, err)
		}
		defer file.Close()

		if _, err := io.Copy(hash, file); err != nil {
			return "", fmt.Errorf("error to read file %s: %w", filePath, err)
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
