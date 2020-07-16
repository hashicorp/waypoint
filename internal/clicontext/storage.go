package clicontext

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/renameio"
)

// Storage is the primary struct for interacting with stored CLI contexts.
// Contexts are always stored directly on disk with one set as the default.
type Storage struct {
	dir string
}

// NewStorage initializes context storage.
func NewStorage(opts ...Option) (*Storage, error) {
	var m Storage
	for _, opt := range opts {
		if err := opt(&m); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

// List lists the contexts that are available.
func (m *Storage) List() ([]string, error) {
	f, err := os.Open(m.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	// Remove all our _-prefixed names which are system settings.
	result := make([]string, 0, len(names))
	for _, n := range names {
		if n[0] == '_' {
			continue
		}

		result = append(result, m.nameFromPath(n))
	}

	return result, nil
}

// Load loads a context with the given name.
func (m *Storage) Load(n string) (*Config, error) {
	return LoadPath(m.configPath(n))
}

// Set will set a new configuration with the given name. This will
// overwrite any existing context of this name.
func (m *Storage) Set(n string, c *Config) error {
	path := m.configPath(n)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		return err
	}

	// If we have no default, set as the default
	def, err := m.Default()
	if err != nil {
		return err
	}
	if def == "" {
		err = m.SetDefault(n)
	}

	return err
}

// Rename renames a context. This will error if the "from" context does not
// exist. If "from" is the default context then the default will be switched
// to "to". If "to" already exists, this will overwrite it.
func (m *Storage) Rename(from, to string) error {
	fromPath := m.configPath(from)
	if _, err := os.Stat(fromPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("context %q does not exist", from)
		}

		return err
	}

	if err := m.Delete(to); err != nil {
		return err
	}

	toPath := m.configPath(to)
	if err := os.Rename(fromPath, toPath); err != nil {
		return err
	}

	def, err := m.Default()
	if err != nil {
		return err
	}
	if def == from {
		return m.SetDefault(to)
	}

	return nil
}

// Delete deletes the context with the given name.
func (m *Storage) Delete(n string) error {
	// Remove it
	err := os.Remove(m.configPath(n))
	if os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return err
	}

	// If our default is this, then unset the default
	def, err := m.Default()
	if err != nil {
		return err
	}
	if def == n {
		err = m.UnsetDefault()
	}

	return err
}

// SetDefault sets the default context to use. If the given context
// doesn't exist, an os.IsNotExist error will be returned.
func (m *Storage) SetDefault(n string) error {
	src := m.configPath(n)
	if _, err := os.Stat(src); err != nil {
		return err
	}

	return renameio.Symlink(src, m.defaultPath())
}

// UnsetDefault unsets the default context.
func (m *Storage) UnsetDefault() error {
	err := os.Remove(m.defaultPath())
	if os.IsNotExist(err) {
		err = nil
	}

	return err
}

// Default returns the name of the default context.
func (m *Storage) Default() (string, error) {
	path, err := os.Readlink(m.defaultPath())
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return "", err
	}

	return m.nameFromPath(path), nil
}

// nameFromPath returns the context name given a path to a context
// HCL file. This is just the name of the file without any extension.
func (m *Storage) nameFromPath(path string) string {
	path = filepath.Base(path)
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			path = path[:i]
			break
		}
	}

	return path
}

func (m *Storage) configPath(n string) string {
	return filepath.Join(m.dir, n+".hcl")
}

func (m *Storage) defaultPath() string {
	return filepath.Join(m.dir, "_default.hcl")
}

type Option func(*Storage) error

// WithDir specifies the directory where context configuration will be stored.
// This doesn't have to exist already but we must have permission to create it.
func WithDir(d string) Option {
	return func(m *Storage) error {
		m.dir = d
		return nil
	}
}
