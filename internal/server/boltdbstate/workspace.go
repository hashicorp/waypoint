package boltdbstate

import (
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

var workspaceBucket = []byte("workspaces")

func init() {
	dbBuckets = append(dbBuckets, workspaceBucket)
	dbIndexers = append(dbIndexers, (*State).workspaceIndexInit)
	schemas = append(schemas, workspaceIndexSchema)
}

// WorkspacePut creates or updates the given Workspace.
//
// Project changes will be ignored
func (s *State) WorkspacePut(workspace *pb.Workspace) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.workspacePut(dbTxn, memTxn, workspace)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// WorkspaceList lists all the workspaces.
func (s *State) WorkspaceList() ([]*pb.Workspace, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(workspaceTableName, workspaceIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	return s.workspaceListFromIter(iter)
}

// WorkspaceListByProject lists all the workspaces used by a project.
func (s *State) WorkspaceListByProject(ref *pb.Ref_Project) ([]*pb.Workspace, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.Get(
		workspaceTableName,
		workspaceProjectIndexName,
		ref.Project,
	)
	if err != nil {
		return nil, err
	}

	return s.workspaceListFromIter(iter)
}

// WorkspaceListByApp lists all the workspaces used by a specific application.
func (s *State) WorkspaceListByApp(ref *pb.Ref_Application) ([]*pb.Workspace, error) {
	// To implement this, we just list by project and filter. Projects
	// don't have that many applications, and the index structure to do this
	// more efficiently would be complicated so its not worth it.
	projectWorkspaces, err := s.WorkspaceListByProject(&pb.Ref_Project{Project: ref.Project})
	if err != nil {
		return nil, err
	}

	// Filter the project workspaces to only include workspaces that also
	// have this application.
	var result []*pb.Workspace
PROJECT_LOOP:
	for _, ws := range projectWorkspaces {
		for _, p := range ws.Projects {
			if strings.ToLower(p.Project.Project) != strings.ToLower(ref.Project) {
				continue
			}

			for _, app := range p.Applications {
				if strings.ToLower(app.Application.Application) != strings.ToLower(ref.Application) {
					continue
				}

				// We have an app match, so add it to our results and
				// reloop on the workspace results
				result = append(result, ws)
				continue PROJECT_LOOP
			}
		}
	}

	return result, nil
}

func (s *State) workspacePut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Workspace,
) error {
	if err := serverptypes.ValidateWorkspace(value); err != nil {
		return err
	}
	prev, err := s.workspaceGet(dbTxn, memTxn, &pb.Ref_Workspace{
		Workspace: value.Name,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		// We ignore NotFound since this function is used to create
		// Workspaces.
		return err
	}
	if err == nil {
		// If we have a previous Workspace, preserve the Projects.
		value.Projects = prev.Projects
	}
	id := []byte(strings.ToLower(value.Name))

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(workspaceBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	_, err = s.workspaceIndexSet(memTxn, id, value)
	return err
}

func (s *State) workspaceListFromIter(iter memdb.ResultIterator) ([]*pb.Workspace, error) {
	var result []*pb.Workspace
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*workspaceIndex)

		var ws *pb.Workspace
		err := s.db.View(func(dbTxn *bolt.Tx) error {
			var err error
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

func (s *State) workspaceGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
) (*pb.Workspace, error) {
	var result pb.Workspace
	b := dbTxn.Bucket(workspaceBucket)
	return &result, dbGet(b, []byte(strings.ToLower(ref.Workspace)), &result)
}

func (s *State) WorkspaceDelete(n string) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.workspaceDelete(dbTxn, memTxn, &pb.Ref_Workspace{Workspace: n})
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

func (s *State) workspaceDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Workspace,
) error {
	w, err := s.workspaceGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return err
	}

	id := s.workspaceId(w)
	if err := dbTxn.Bucket(workspaceBucket).Delete([]byte(id)); err != nil {
		return err
	}

	// Delete from memdb
	if _, err := memTxn.DeleteAll(workspaceTableName, workspaceIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}

	return nil
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
	tsProto := timestamppb.New(ts)
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

// workspaceListProjects lists the project information for all the workspaces
// that a specific project is in.
func (s *State) workspaceListProjects(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) ([]*pb.Workspace_Project, error) {
	iter, err := memTxn.Get(
		workspaceTableName,
		workspaceProjectIndexName,
		ref.Project,
	)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workspace_Project
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*workspaceIndex)

		ws, err := s.workspaceFromDB(dbTxn, idx.Id)
		if err != nil {
			return nil, err
		}
		for _, p := range ws.Projects {
			if strings.ToLower(p.Project.Project) == strings.ToLower(ref.Project) {
				// This gets set only for this API call (as documented in the proto)
				p.Workspace = &pb.Ref_Workspace{Workspace: ws.Name}

				result = append(result, p)
				break
			}
		}
	}

	return result, nil
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

// workspaceInitApp finds the given app or creates a new one if it doesn't exist.
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

func (s *State) workspaceId(w *pb.Workspace) []byte {
	return []byte(strings.ToLower(w.Name))
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

	// Index all the projects that are in this workspace
	for _, p := range wspb.Projects {
		rec.Projects = append(rec.Projects, p.Project.Project)
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

			workspaceProjectIndexName: {
				Name:         workspaceProjectIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.StringSliceFieldIndex{
					Field:     "Projects",
					Lowercase: true,
				},
			},
		},
	}
}

const (
	workspaceTableName        = "workspace-index"
	workspaceIdIndexName      = "id"
	workspaceProjectIndexName = "project"
)

type workspaceIndex struct {
	Id       string   // Id is the name of the workspace lowercased
	Projects []string // Projects that are part of this workspace
}
