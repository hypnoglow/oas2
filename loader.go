package oas

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

// LoadSpecOption is an option for LoadSpec.
type LoadSpecOption struct {
	optionType optionType
	value      interface{}
}

// CacheDir is an option for LoadSpec that sets cache directory for faster
// spec loads.
func CacheDir(dir string) *LoadSpecOption {
	return &LoadSpecOption{
		optionType: cacheDirOption,
		value:      dir,
	}
}

// LoadSpec opens an OpenAPI Specification v2.0 document, expands all references within it,
// then validates the spec and returns spec document.
func LoadSpec(fpath string, opts ...*LoadSpecOption) (document *loads.Document, err error) {
	var cacheDir string

	for _, opt := range opts {
		if val, ok := opt.value.(string); ok && opt.optionType == cacheDirOption {
			cacheDir = val
		}
	}

	// Load from cache.

	var cacheFilename string
	var hashSum string
	if cacheDir != "" {
		hashSum, err = hashFile(fpath)
		if err != nil {
			return nil, errors.Wrap(err, "calculate file hash")
		}

		cacheFilename = filepath.Join(cacheDir, fmt.Sprintf("spec-%s", hashSum))
		document, err = loads.JSONSpec(cacheFilename)
		if err == nil {
			return document, nil
		}

		// if an error occurred - ignore it and continue.
	}

	// Load regularly.

	document, err = loads.Spec(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load spec")
	}

	document, err = document.Expanded(&spec.ExpandOptions{RelativeBase: fpath})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to expand spec")
	}

	if err = validate.Spec(document, strfmt.Default); err != nil {
		return nil, errors.Wrap(err, "Spec is invalid")
	}

	if cacheDir == "" {
		return document, nil
	}

	// Cache expanded spec.

	if err = os.MkdirAll(filepath.Dir(cacheFilename), 0700); err != nil {
		return document, errors.Wrap(err, "create cache dir")
	}

	f, err := os.Create(cacheFilename)
	if err != nil {
		return document, errors.Wrap(err, "create cache file")
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(document.Spec()); err != nil {
		return document, errors.Wrap(err, "write cache file")
	}

	return document, nil
}

type optionType string

const (
	cacheDirOption optionType = "cache directory"
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
