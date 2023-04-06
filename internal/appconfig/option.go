// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfig

import (
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint/internal/plugin"
)

// Option is used to configure NewWatcher.
type Option func(w *Watcher) error

// WithLogger sets the logger for the Watcher. If no logger is specified,
// the Watcher will use the default logger (hclog.L() value).
func WithLogger(log hclog.Logger) Option {
	return func(w *Watcher) error {
		if log != nil {
			w.log = log
		}

		return nil
	}
}

// WithPlugins sets a map of already-launched plugins to use for dynamic
// configuration sources.
func WithPlugins(ps map[string]*plugin.Instance) Option {
	return func(w *Watcher) error {
		w.plugins = ps
		return nil
	}
}

// Contains the information about a file that should be written to the application's
// current directory.
type FileContent struct {
	Path string
	Data []byte
}

// UpdatedConfig contains any updated configuration that needs to be applied to the
// application.
type UpdatedConfig struct {
	// Indicates that EnvVars is what the application should be using. This is an
	// explicit flag because EnvVars might be reset to nil, meaning the application
	// should remove all its configuration
	UpdatedEnv bool

	// This is the list of env vars in key=value format that the application should
	// know about.
	EnvVars []string

	// DeletedEnvVars will be a list of keys that were unset during this change
	// that was previously set. Note that if the watcher is running with
	// the `WithOriginalEnv` configuration, then an unset config MIGHT not be
	// here if it is replaced with the original value in EnvVars.
	DeletedEnvVars []string

	// Indicates that Files is what should be presented on disk. This is an explicit
	// flag to match UpdatedEnv.
	UpdatedFiles bool

	// Files is the list of file paths and contents that the should be on disk for the
	// application to read.
	Files []*FileContent
}

// WithNotify notifies a channel whenever there are changes to the
// configuration values. This will stop receiving values when the watcher
// is closed.
//
// Updates will block when attempting to send on this channel. However,
// while blocking, multiple updates may occur that will be coalesced to a
// follow up update when the channel send succeeds. Therefore, receivers
// will always eventually receive the full current env list, but may miss
// intermediate sets if they are slow to receive.
func WithNotify(ch chan<- *UpdatedConfig) Option {
	return func(w *Watcher) error {
		// Start the goroutine for watching. If there is an error during
		// init, NewWatcher calls Close so these will be cleaned up.
		go w.notify(w.bgCtx, ch)
		return nil
	}
}

// WithRefreshInterval sets the interval between checking for new values
// from config source plugins that don't support edge triggers.
//
// NOTE(mitchellh): At the time of writing, we don't support edge triggered
// plugins at all, but we plan to at some point so the docs reflect that.
func WithRefreshInterval(d time.Duration) Option {
	return func(w *Watcher) error {
		w.refreshInterval = d
		return nil
	}
}

// WithDynamicEnabled sets whether we allow dynamic sources or not.
// This defaults to true.
//
// If this is disabled, then all dynamic config requests are ignored.
// They aren't set to empty values or anything, they simply aren't set
// at all.
func WithDynamicEnabled(v bool) Option {
	return func(w *Watcher) error {
		w.dynamicEnabled = v
		return nil
	}
}

// WithOriginalEnv sets an "original" set of environment variables. When
// an environment variable value is unset, it will use a value from here
// if the key exists.
//
// The value of v is the same as os.Environ().
func WithOriginalEnv(v []string) Option {
	return func(w *Watcher) error {
		w.originalEnv = v
		return nil
	}
}
