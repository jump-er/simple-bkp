package storage

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func FileHashes(path string) (md5sum string, sha256sum string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	hMD5 := md5.New()
	hSHA := sha256.New()

	_, err = io.Copy(io.MultiWriter(hMD5, hSHA), f)
	if err != nil {
		return "", "", err
	}

	md5sum = fmt.Sprintf("%x", hMD5.Sum(nil))
	sha256sum = fmt.Sprintf("%x", hSHA.Sum(nil))
	return
}
