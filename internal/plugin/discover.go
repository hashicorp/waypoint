package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"

	"github.com/hashicorp/waypoint/internal/config"
)

// Discover finds the given plugin and returns the command for it. The command
// can subsequently be used with Factory to build a factory for a specific
// plugin type. If the plugin is not found `(nil, nil)` is returned.
//
// The plugin binary must have the form "waypoint-plugin-<name>" (with a
// ".exe" extension on Windows).
//
// This will search the paths given. You can use DefaultPaths() to get
// the default set of paths.
//
func Discover(cfg *config.Plugin, paths []string) (*exec.Cmd, error) {
	// Expected filename
	expected := "waypoint-plugin-" + cfg.Name
	if runtime.GOOS == "windows " {
		expected += ".exe"
	}

	// Search our paths
	for _, path := range paths {
		path = filepath.Join(path, expected)

		_, err := os.Stat(path)
		if err == nil {
			cmd := exec.Command(path)
			return cmd, nil
		}

		if os.IsNotExist(err) {
			continue
		}

		return nil, err
	}

	return nil, nil
}

// DefaultPaths returns the default search paths for plugins. These are:
//
//   * pwd given
//   * "$pwd/.waypoint/plugins"
//   * "$XDG_CONFIG_DIR/waypoint/plugins"
//
func DefaultPaths(pwd string) ([]string, error) {
	xdgPath, err := xdg.ConfigFile("waypoint/plugins/.ignore")
	if err != nil {
		return nil, err
	}

	return []string{
		pwd,
		filepath.Join(pwd, ".waypoint", "plugins"),
		filepath.Dir(xdgPath),
	}, nil
}
