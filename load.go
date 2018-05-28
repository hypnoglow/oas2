package oas

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/pkg/errors"
)

// LoadOptions represent options that are used on specification load.
type LoadOptions struct {
	host       string
	schemes    []string
	appVersion string

	cacheDir string
}

// LoadOption is option to use when loading specification.
type LoadOption func(*LoadOptions)

// LoadSetHost returns option that sets specification host.
func LoadSetHost(host string) LoadOption {
	return func(o *LoadOptions) {
		o.host = host
	}
}

// LoadSetSchemes returns option that sets specification schemes.
func LoadSetSchemes(schemes []string) LoadOption {
	return func(o *LoadOptions) {
		o.schemes = schemes
	}
}

// LoadSetAPIVersion returns option that sets application API version.
func LoadSetAPIVersion(version string) LoadOption {
	return func(o *LoadOptions) {
		o.appVersion = version
	}
}

// LoadCacheDir returns option that allows to load expanded spec from cache.
func LoadCacheDir(dir string) LoadOption {
	return func(o *LoadOptions) {
		o.cacheDir = dir
	}
}

// LoadFile loads OpenAPI specification from file.
func LoadFile(fpath string, opts ...LoadOption) (*loads.Document, error) {
	options := LoadOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	document, err := loadDocument(fpath, options.cacheDir)
	if err != nil {
		return nil, err
	}

	if options.host != "" {
		document.Spec().Host = options.host
		document.OrigSpec().Host = options.host
	}

	if options.schemes != nil {
		document.Spec().Schemes = options.schemes
		document.OrigSpec().Schemes = options.schemes
	}

	if options.appVersion != "" {
		document.Spec().Info.Version = options.appVersion
		document.OrigSpec().Info.Version = options.appVersion
	}

	return document, nil
}

func loadDocument(fpath, cacheDir string) (*loads.Document, error) {
	document, err := loads.Spec(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "load spec from file")
	}

	hashSum, err := hashFile(fpath)
	if err != nil {
		return nil, errors.Wrap(err, "calculate file hash")
	}

	if exp, err := loadExpandedFromCache(cacheDir, hashSum); err == nil {
		// When document loaded from cache, it is safe to use exp.Raw()
		doc, e := loads.Embedded(document.Raw(), exp.Raw())
		return doc, errors.Wrap(e, "create embedded document")
	}

	// If cannot load from cache for some reason - expand original spec.

	// We assume that everything cached is valid, but when cache is empty -
	// we need to validate the original document.
	if err = validate.Spec(document, strfmt.Default); err != nil {
		return nil, errors.Wrap(err, "validate spec")
	}

	exp, err := document.Expanded(&spec.ExpandOptions{RelativeBase: fpath})
	if err != nil {
		return nil, errors.Wrap(err, "expand spec")
	}

	if err = saveExpandedToCache(exp, cacheDir, hashSum); err != nil {
		return nil, errors.Wrap(err, "save expanded spec to cache")
	}

	// To use expanded document right away, we need to get raw from it.
	// WARNING: When document is expanded in memory like above, exp.Raw() still
	// returns not expanded spec, so do not try to use it here.
	expBytes, err := exp.Spec().MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "convert expanded spec to raw")
	}

	doc, err := loads.Embedded(document.Raw(), json.RawMessage(expBytes))
	return doc, errors.Wrap(err, "create embedded document")
}

// loadExpandedFromCache loads OpenAPI document from cache if cacheDir is not empty.
func loadExpandedFromCache(cacheDir, fpath string) (*loads.Document, error) {
	if cacheDir == "" {
		return nil, errors.New("cache dir is empty")
	}

	cacheFilename := filepath.Join(cacheDir, fpath) + ".json"

	return loads.JSONSpec(cacheFilename)
}

// saveExpandedToCache saves OpenAPI document to cache if cacheDir is not empty.
func saveExpandedToCache(expandedDoc *loads.Document, cacheDir, fpath string) error {
	if cacheDir == "" {
		return nil
	}

	cacheFilename := filepath.Join(cacheDir, fpath) + ".json"

	if err := os.MkdirAll(filepath.Dir(cacheFilename), 0700); err != nil {
		return errors.Wrap(err, "create cache dir")
	}

	f, err := os.Create(cacheFilename)
	if err != nil {
		return errors.Wrap(err, "create cache file")
	}
	defer f.Close()

	if err = json.NewEncoder(f).Encode(expandedDoc.Spec()); err != nil {
		return errors.Wrap(err, "write cache file")
	}
	return nil
}
