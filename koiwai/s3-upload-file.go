package koiwai

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

func (s *S3Service) S3UploadFile(filename string) *[]byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer f.Close()

	h := sha256.New()

	if _, err := io.Copy(h, f); err != nil {
		log.Println(err)
		return nil
	}

	hashBytes := h.Sum(nil)
	hashAlreadyExists, err := s.KoiwaiDb.NewSelect().
		Model(&Media{}).
		Where("hash = ?", hashBytes).
		Exists(context.Background())

	if hashAlreadyExists {
		return &hashBytes
	}

	hashString := base64.URLEncoding.EncodeToString(hashBytes)

	_, err = s.S3Client.FPutObject(context.Background(), s.BucketName, hashString, filename, minio.PutObjectOptions{})

	if err != nil {
		log.Println(err)
		return nil
	}

	if _, err := s.KoiwaiDb.NewInsert().Model(&Media{Hash: hashBytes}).Exec(context.Background()); err != nil {
		return nil
	}

	return &hashBytes
}
