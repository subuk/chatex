package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func GetSha1Hex(src io.ReadSeeker) string {
	hasher := sha1.New()
	io.Copy(hasher, src)
	src.Seek(0, 0)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetUploadFileSize(f io.Seeker) int64 {
	f.Seek(0, 0)
	size, err := f.Seek(0, 2)
	f.Seek(0, 0)

	if err != nil {
		panic(err)
	}
	return size
}
