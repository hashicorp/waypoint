package singleprocess

import (
	"github.com/hashicorp/go-memdb"
)

// memdbSchema is the schema to setup for the in-memory database. This
// can be populated with init functions in the respective files.
var memdbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{},
}

func memdbInit() (*memdb.MemDB, error) {
	return memdb.NewMemDB(memdbSchema)
}
