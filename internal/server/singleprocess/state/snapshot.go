package state

import (
	"io"

	"github.com/boltdb/bolt"
)

// CreateSnapshot creates a database snapshot and writes it to the given writer.
func (s *State) CreateSnapshot(w io.Writer) error {
	return s.db.View(func(dbTxn *bolt.Tx) error {
		_, err := dbTxn.WriteTo(w)
		return err
	})
}
