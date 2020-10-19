package state

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func init() {
	schemas = append(schemas, workspaceIndexSchema)
}

// WorkspaceList lists all the workspaces.
func (s *State) WorkspaceList() ([]*pb.Workspace, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(
		workspaceIndexTableName,
		workspaceIndexIdIndexName+"_prefix",
		"",
	)
	if err != nil {
		return nil, err
	}

	wsMap := map[string]*pb.Workspace{}
	appMap := map[string]*pb.Workspace_Application{}
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}

		idx := raw.(*workspaceIndexRecord)

		// Get our workspace
		ws, ok := wsMap[idx.Name]
		if !ok {
			ws = &pb.Workspace{Name: idx.Name}
			wsMap[idx.Name] = ws
		}

		// Get our application
		key := fmt.Sprintf("%s/%s/%s",
			idx.Name,
			idx.Project,
			idx.App,
		)

		// If we don't have the app yet, create it
		app, ok := appMap[key]
		if !ok {
			app = &pb.Workspace_Application{
				Application: &pb.Ref_Application{
					Project:     idx.Project,
					Application: idx.App,
				},
			}

			// Append our apps
			ws.Applications = append(ws.Applications, app)

			// Keep track of it so we can reuse this later
			appMap[key] = app
		}

		// Get our current timestamp to compare to this record. We can
		// make this more efficient by storing this but its not a big deal
		// right now.
		var currentActive time.Time
		if app.ActiveTime != nil {
			currentActive, err = ptypes.Timestamp(app.ActiveTime)
			if err != nil {
				return nil, err
			}
		}

		// If this time is later than our current active, we store it.
		if idx.LastActiveAt.After(currentActive) {
			// Set our new current
			app.ActiveTime, err = ptypes.TimestampProto(idx.LastActiveAt)
			if err != nil {
				return nil, err
			}

			// Get our current workspace time. Similiar to above, can
			// optimize later by caching this.
			var currentMax time.Time
			if ws.ActiveTime != nil {
				currentMax, err = ptypes.Timestamp(ws.ActiveTime)
				if err != nil {
					return nil, err
				}
			}

			// Set our max value
			if currentMax.IsZero() || idx.LastActiveAt.After(currentMax) {
				ws.ActiveTime = app.ActiveTime
			}
		}
	}

	result := make([]*pb.Workspace, 0, len(wsMap))
	for _, v := range wsMap {
		result = append(result, v)
	}

	return result, nil
}

// WorkspaceGet gets a workspace with a specific name. If it doesn't exist,
// this will return an error with codes.NotFound.
func (s *State) WorkspaceGet(n string) (*pb.Workspace, error) {
	// We implement this in terms of list for now.
	wsList, err := s.WorkspaceList()
	if err != nil {
		return nil, err
	}

	for _, ws := range wsList {
		if strings.EqualFold(ws.Name, n) {
			return ws, nil
		}
	}

	return nil, status.Errorf(codes.NotFound,
		"not found for name: %q", n)
}

// workspaceTouch creates an initial record for the given workspace or
// returns one if it already exists. It also updates the LastActiveAt time
// for this resource.
func (s *State) workspaceTouch(
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
	app *pb.Ref_Application,
	resource string,
	ts time.Time,
) (*workspaceIndexRecord, error) {
	rec, err := s.workspaceGet(memTxn, ref, app, resource)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		rec = &workspaceIndexRecord{
			Name:     ref.Workspace,
			Project:  app.Project,
			App:      app.Application,
			Resource: resource,
		}
	}

	if v := rec.LastActiveAt; v.IsZero() || ts.After(v) {
		// Set the new last active at
		rec.LastActiveAt = ts
	}

	// Store
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
	resource string,
) (*workspaceIndexRecord, error) {
	raw, err := memTxn.First(
		workspaceIndexTableName,
		workspaceIndexIdIndexName,
		ref.Workspace,
		app.Project,
		app.Application,
		resource,
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
			workspaceIndexIdIndexName: {
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

						&memdb.StringFieldIndex{
							Field:     "Resource",
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

	// Resource is the resource that is being managed within this workspace
	// for a given application. This is chosen by the resource itself but
	// should be unique. For example, app operations such as build use the
	// lowercase struct name.
	Resource string

	// LastActiveAt is a timestamp when a resource of this type was last
	// "active". The definition of active is typically created or modified
	// but its up to the resource type.
	LastActiveAt time.Time
}
