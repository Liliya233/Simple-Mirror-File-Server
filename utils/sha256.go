package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// 获取文件Hash
func GetFileHash(fname string) string {
	file, err := os.Open(fname)
	if err != nil {
		return ""
	}
	defer file.Close()
	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return ""
	}
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}
