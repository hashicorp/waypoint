package datadir

import (
	"path/filepath"
)

// App is an implementation of Dir that encapsulates the directories for a
// single app.
type App struct {
	Dir
}

// Component returns a Dir implementation scoped to a specific component.
func (d *App) Component(typ, name string) (*Component, error) {
	dir, err := NewScopedDir(d, filepath.Join("component", typ, name))
	if err != nil {
		return nil, err
	}

	return &Component{Dir: dir}, nil
}
