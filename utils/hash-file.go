package utils

import (
	"io"
	"log"
	"os"

	"github.com/minio/sha256-simd"
)

func HashFile(filename string) *[]byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return nil
	}

	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Printf("Error hashing file: %s\n", err)
		return nil
	}

	hash := h.Sum(nil)
	return &hash
}
