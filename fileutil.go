package oas

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"

	"github.com/pkg/errors"
)

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "open spec file")
	}
	defer f.Close()

	h := sha512.New512_256()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "copy file to hash")
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
