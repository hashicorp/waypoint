package state

import (
	"github.com/hashicorp/go-memdb"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func init() {
	schemas = append(schemas, workspaceIndexSchema)
}

// workspaceInit creates an initial record for the given workspace or
// returns one if it already exists.
func (s *State) workspaceInit(
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
	app *pb.Ref_Application,
) (*workspaceIndexRecord, error) {
	rec, err := s.workspaceGet(memTxn, ref, app)
	if err != nil {
		return nil, err
	}
	if rec != nil {
		return rec, nil
	}

	rec = &workspaceIndexRecord{
		Name:    ref.Workspace,
		Project: app.Project,
		App:     app.Application,
	}
	return rec, s.workspacePut(memTxn, rec)
}

// workspacePut updates the workspace record.
func (s *State) workspacePut(
	memTxn *memdb.Txn,
	rec *workspaceIndexRecord,
) error {
	return memTxn.Insert(workspaceIndexTableName, rec)
}

func (s *State) workspaceGet(
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
	app *pb.Ref_Application,
) (*workspaceIndexRecord, error) {
	raw, err := memTxn.First(
		workspaceIndexTableName,
		workspaceIndexIdIndexName,
		ref.Workspace,
		app.Project,
		app.Application,
	)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	return raw.(*workspaceIndexRecord), nil
}

func workspaceIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: workspaceIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			workspaceIndexIdIndexName: &memdb.IndexSchema{
				Name:         workspaceIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Name",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},
					},
				},
			},
		},
	}
}

const (
	workspaceIndexTableName   = "workspace-index"
	workspaceIndexIdIndexName = "id"
)

type workspaceIndexRecord struct {
	Name    string
	Project string
	App     string
}
