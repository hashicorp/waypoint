package state

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gofrs/flock"
	"github.com/hashicorp/go-hclog"
	"github.com/natefinch/atomic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateSnapshot creates a database snapshot and writes it to the given writer.
func (s *State) CreateSnapshot(w io.Writer) error {
	return s.db.View(func(dbTxn *bolt.Tx) error {
		_, err := dbTxn.WriteTo(w)
		return err
	})
}

// StageRestoreSnapshot stages a database restore for the next server restart.
// This will create a temporary file alongside the data file so we must have
// write access to the directory containing the database.
func (s *State) StageRestoreSnapshot(r io.Reader) error {
	log := s.log.Named("restore")
	log.Warn("beginning to stage snapshot restore")

	ri := newRestoreInfo(log, s.db)
	if err := ri.Lock(); err != nil {
		return err
	}
	defer ri.Unlock()

	// Get our file info
	fi, err := os.Stat(s.db.Path())
	if err != nil {
		return err
	}

	// Open our temporary file and copy the restore contents into it.
	log.Info("copying the snapshot data to a temporary path", "path", ri.StageTempPath)
	tempF, err := os.OpenFile(ri.StageTempPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		log.Error("error creating temporary path", "err", err)
		return err
	}
	_, err = io.Copy(tempF, r)
	tempF.Close()
	if err != nil {
		log.Error("error copying snapshot data", "err", err)
		return err
	}

	// Replace our file
	log.Info("atomically replacing file", "src", ri.StageTempPath, "dest", ri.StagePath)
	if err := atomic.ReplaceFile(ri.StageTempPath, ri.StagePath); err != nil {
		log.Error("error replacing file", "err", err)
		return err
	}

	// Open the new database
	log.Warn("snapshot staged for restore on next startup")
	return nil
}

// finalizeRestore checks for any staged restore and moves it into place.
// This will error if it fails for any reason which may prevent startup
// but we have to prevent startup because the user wanted a restore.
func finalizeRestore(log hclog.Logger, db *bolt.DB) (*bolt.DB, error) {
	log.Debug("checking if DB restore is requested")
	ri := newRestoreInfo(log, db)
	if err := ri.Lock(); err != nil {
		return db, err
	}
	defer ri.Unlock()

	_, err := os.Stat(ri.StagePath)
	if os.IsNotExist(err) {
		log.Debug("no restore file found, no DB restore requested")
		return db, nil
	}
	if err != nil {
		log.Error("error checking for restore file", "err", err)
		return db, err
	}

	log.Warn("restore file found, will initiate database restore")

	// Close our DB, we will reopen with the new one
	if err := db.Close(); err != nil {
		log.Error("failed to close db for restore", "err", err)
		return db, err
	}

	// Get our file info
	var mode os.FileMode = 0666
	if fi, err := os.Stat(ri.DBPath); err == nil {
		mode = fi.Mode()
	} else if !os.IsNotExist(err) {
		return db, err
	}

	// Replace our file
	log.Info("atomically replacing db file", "src", ri.StagePath, "dest", ri.DBPath)
	if err := atomic.ReplaceFile(ri.StagePath, ri.DBPath); err != nil {
		log.Error("error replacing file", "err", err)
		return db, err
	}

	// Reopen the DB
	log.Info("reopening database", "path", ri.DBPath)
	db, err = bolt.Open(ri.DBPath, mode, &bolt.Options{
		Timeout: 2 * time.Second,
	})
	if err != nil {
		log.Error("error reopening db", "err", err)
		return db, err
	}

	log.Warn("database restore successful")
	return db, nil
}

type restoreInfo struct {
	// DBPath is the final database path.
	DBPath string

	// StagePath is the final path where the staged restore data should
	// be placed. If this path exists, the data is expected to be valid
	// and not corrupted.
	StagePath string

	// StageTempPath is where the staged restore data should be temporarily
	// written to while it is still being loaded. This data if it exists
	// may be corrupted or incomplete until it is atomically moved to
	// StagePath.
	StageTempPath string

	fl  *flock.Flock
	log hclog.Logger
}

// newRestoreInfo gets the restore info from the given DB.
func newRestoreInfo(log hclog.Logger, db *bolt.DB) *restoreInfo {
	// Get our current directory
	destPath := db.Path()
	dir := filepath.Dir(destPath)

	// Paths to our restore file
	stagePath := filepath.Join(dir, "waypoint-restore.db")
	tempPath := stagePath + ".temp"
	lockPath := stagePath + ".lock"

	return &restoreInfo{
		DBPath:        destPath,
		StagePath:     stagePath,
		StageTempPath: tempPath,
		fl:            flock.New(lockPath),
		log:           log,
	}
}

// Lock locks the restore lockfile or returns an error if this failed.
// If the return value is nil, then Unlock must be called to unlock.
func (r *restoreInfo) Lock() error {
	// Create a file lock to ensure only one restore is happening at a time
	r.log.Trace("acquiring file lock for restore", "path", r.fl.String())
	locked, err := r.fl.TryLock()
	if err != nil {
		r.log.Error("error acquiring file lock", "err", err)
		return err
	}
	if !locked {
		r.log.Error("error acquiring file lock, lock already held")
		return status.Errorf(codes.Aborted,
			"failed to acquire file lock for restore, another restore may already be active")
	}

	return nil
}

// Unlock unlocks the file lock. This is only safe to call if a lock
// was successfully acquired.
func (r *restoreInfo) Unlock() error {
	r.log.Trace("releasing file lock for restore", "path", r.fl.String())
	return r.fl.Unlock()
}
