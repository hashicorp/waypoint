package datadir

import (
	"os"
	"path/filepath"
)

// TODO(mitchellh): we use deeply nested directories here which isn't
// going to work on Windows (due to MAX_PATH). We should have an alternate
// implementation for Windows.

// TODO(mitchellh): tests! like any tests

// Dir is the interface implemented so that consumers can store data
// locally in a consistent way.
type Dir interface {
	// CacheDir returns the path to a folder that can be used for
	// cache data. This directory may not be empty if a previous run
	// stored data, but it may also be emptied at any time between runs.
	CacheDir() string

	// DataDir returns the path to a folder that can be used for data
	// that is persisted between runs.
	DataDir() string
}

// basicDir implements Dir in the simplest possible way.
type basicDir struct {
	cacheDir string
	dataDir  string
}

// CacheDir impl Dir
func (d *basicDir) CacheDir() string { return d.cacheDir }

// DataDir impl Dir
func (d *basicDir) DataDir() string { return d.dataDir }

// newRootDir creates a basicDir for the root directory which puts
// data at <path>/cache, etc.
func newRootDir(path string) (Dir, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(path, "cache")
	dataDir := filepath.Join(path, "data")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &basicDir{cacheDir: cacheDir, dataDir: dataDir}, nil
}

// NewBasicDir creates a Dir implementation with a manually specified
// set of directories.
func NewBasicDir(cacheDir, dataDir string) Dir {
	return &basicDir{cacheDir: cacheDir, dataDir: dataDir}
}

// NewScopedDir creates a ScopedDir for the given parent at the relative
// child path of path. The caller should take care that multiple scoped
// dirs with overlapping paths are not created, since they could still
// collide.
func NewScopedDir(parent Dir, path string) (Dir, error) {
	cacheDir := filepath.Join(parent.CacheDir(), path)
	dataDir := filepath.Join(parent.DataDir(), path)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	return &basicDir{cacheDir: cacheDir, dataDir: dataDir}, nil
}

var _ Dir = (*basicDir)(nil)
