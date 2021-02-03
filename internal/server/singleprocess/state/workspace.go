package state

import (
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var (
	workspaceBucket = []byte("workspaces")
)

func init() {
	dbBuckets = append(dbBuckets, workspaceBucket)
	dbIndexers = append(dbIndexers, (*State).workspaceIndexInit)
	schemas = append(schemas, workspaceIndexSchema)
}

// WorkspaceList lists all the workspaces.
func (s *State) WorkspaceList() ([]*pb.Workspace, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(workspaceTableName, workspaceIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Workspace
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*workspaceIndex)

		var ws *pb.Workspace
		err := s.db.View(func(dbTxn *bolt.Tx) error {
			ws, err = s.workspaceFromDB(dbTxn, idx.Id)
			return err
		})
		if err != nil {
			return nil, err
		}

		result = append(result, ws)
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

// workspaceTouchApp updates a workspace with the given project/application
// usage. This will create the workspace if it does not exist. This updates
// the LastActiveAt times.
func (s *State) workspaceTouchApp(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
	app *pb.Ref_Application,
	ts time.Time,
) error {
	id := strings.ToLower(ref.Workspace)

	// Look up the workspace
	raw, err := memTxn.First(workspaceTableName, workspaceIdIndexName, id)
	if err != nil {
		return err
	}

	var ws *pb.Workspace
	if raw != nil {
		// If we have a previous record, load it.
		ws, err = s.workspaceFromDB(dbTxn, id)
		if err != nil {
			return err
		}
	}
	if ws == nil {
		// If we have no workspace, create a new one.
		ws = &pb.Workspace{Name: ref.Workspace}
	}

	// Initialize the project
	wsP, err := s.workspaceInitProject(ws, &pb.Ref_Project{
		Project: app.Project,
	})
	if err != nil {
		return err
	}

	// Initialize the app
	wsApp, err := s.workspaceInitApp(wsP, app)
	if err != nil {
		return err
	}

	// Update our timestamps
	tsProto, err := ptypes.TimestampProto(ts)
	if err != nil {
		return err
	}
	ws.ActiveTime = tsProto
	wsP.ActiveTime = tsProto
	wsApp.ActiveTime = tsProto

	// Store and update index
	if err := dbPut(dbTxn.Bucket(workspaceBucket), []byte(id), ws); err != nil {
		return err
	}
	_, err = s.workspaceIndexSet(memTxn, []byte(id), ws)
	return err
}

// workspaceUpdateProjectDataRef updates the latest data ref used for a project.
func (s *State) workspaceUpdateProjectDataRef(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
	project *pb.Ref_Project,
	dataRef *pb.Job_DataSource_Ref,
) error {
	id := strings.ToLower(ref.Workspace)

	// Look up the workspace
	raw, err := memTxn.First(workspaceTableName, workspaceIdIndexName, id)
	if err != nil {
		return err
	}

	var ws *pb.Workspace
	if raw != nil {
		// If we have a previous record, load it.
		ws, err = s.workspaceFromDB(dbTxn, id)
		if err != nil {
			return err
		}
	}
	if ws == nil {
		// If we have no workspace, create a new one.
		ws = &pb.Workspace{Name: ref.Workspace}
	}

	// Initialize the project
	wsP, err := s.workspaceInitProject(ws, project)
	if err != nil {
		return err
	}

	// Set the project data ref
	wsP.DataSourceRef = dataRef

	// Store and update index
	if err := dbPut(dbTxn.Bucket(workspaceBucket), []byte(id), ws); err != nil {
		return err
	}
	_, err = s.workspaceIndexSet(memTxn, []byte(id), ws)
	return err
}

// workspaceFromDB loads the Workspace structure from disk.
func (s *State) workspaceFromDB(dbTxn *bolt.Tx, id string) (*pb.Workspace, error) {
	var result pb.Workspace
	b := dbTxn.Bucket(workspaceBucket)
	return &result, dbGet(b, []byte(strings.ToLower(id)), &result)
}

// workspaceInitProject finds the given project or creates a new one on the
// workspace if it doesn't exist. This does not persist to any database.
func (s *State) workspaceInitProject(
	ws *pb.Workspace,
	ref *pb.Ref_Project,
) (*pb.Workspace_Project, error) {
	// Search for an existing project
	for _, p := range ws.Projects {
		if strings.EqualFold(p.Project.Project, ref.Project) {
			return p, nil
		}
	}

	// If we didn't find one, then create it
	p := &pb.Workspace_Project{Project: ref}
	ws.Projects = append(ws.Projects, p)
	return p, nil
}

// workspaceInitApp finds the given app or creates a new one if it doensn't exist.
// This does not persist to any database.
func (s *State) workspaceInitApp(
	wsProject *pb.Workspace_Project,
	ref *pb.Ref_Application,
) (*pb.Workspace_Application, error) {
	// Basic validation to avoid data corruption: the projects must match
	if !strings.EqualFold(wsProject.Project.Project, ref.Project) {
		return nil, status.Errorf(codes.Internal,
			"application project must match workspace project")
	}

	// Search for an existing project
	for _, app := range wsProject.Applications {
		if strings.EqualFold(app.Application.Application, ref.Application) {
			return app, nil
		}
	}

	// If we didn't find one, then create it
	app := &pb.Workspace_Application{Application: ref}
	wsProject.Applications = append(wsProject.Applications, app)
	return app, nil
}

// workspaceIndexInit initializes the config index from persisted data.
func (s *State) workspaceIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(workspaceBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Workspace
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}

		_, err := s.workspaceIndexSet(memTxn, k, &value)
		if err != nil {
			return err
		}

		return nil
	})
}

// workspaceIndexSet writes an index record for a single workspace.
func (s *State) workspaceIndexSet(
	txn *memdb.Txn,
	id []byte,
	wspb *pb.Workspace,
) (*workspaceIndex, error) {
	rec := &workspaceIndex{
		Id: wspb.Name,
	}

	// Insert the index
	return rec, txn.Insert(workspaceTableName, rec)
}

func workspaceIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: workspaceTableName,
		Indexes: map[string]*memdb.IndexSchema{
			workspaceIdIndexName: {
				Name:         workspaceIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	workspaceTableName   = "workspace-index"
	workspaceIdIndexName = "id"
)

type workspaceIndex struct {
	Id string // Id is the name of the workspace lowercased
}
