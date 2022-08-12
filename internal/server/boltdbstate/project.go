package boltdbstate

import (
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

var projectBucket = []byte("project")

func init() {
	dbBuckets = append(dbBuckets, projectBucket)
	dbIndexers = append(dbIndexers, (*State).projectIndexInit)
	schemas = append(schemas, projectIndexSchema)
}

// ProjectPut creates or updates the given project.
//
// Application changes will be ignored, you must use the Application APIs.
func (s *State) ProjectPut(p *pb.Project) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		prev, err := s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
			Project: p.Name,
		})
		if err != nil && status.Code(err) != codes.NotFound {
			// We ignore NotFound since this function is used to create projects.
			return err
		}
		if err == nil {
			// If we have a previous project, preserve the applications.
			p.Applications = prev.Applications
		}

		return s.projectPut(dbTxn, memTxn, p)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ProjectGet gets a project by reference.
func (s *State) ProjectGet(ref *pb.Ref_Project) (*pb.Project, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var result *pb.Project
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.projectGet(dbTxn, memTxn, ref)
		return err
	})

	return result, err
}

// ProjectDelete deletes a project by reference. This is a complete data
// delete. This will delete all operations associated with this project
// as well.
func (s *State) ProjectDelete(ref *pb.Ref_Project) error {
	// Get our project to delete
	project, err := s.ProjectGet(ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return err
	}

	// We perform all of our reads before our write to avoid the deadlocked state
	var builds []*pb.Build
	var artifacts []*pb.PushedArtifact
	var deployments []*pb.Deployment
	var releases []*pb.Release
	var statusReports []*pb.StatusReport
	var workspaces []*pb.Workspace_Project
	var triggers []*pb.Trigger
	var pipelines []*pb.Pipeline
	if err = s.db.View(func(dbTxn *bolt.Tx) error {
		for _, app := range project.Applications {
			appRef := &pb.Ref_Application{
				Application: app.Name,
				Project:     project.Name,
			}
			if builds, err = s.BuildList(appRef); err != nil {
				return err
			}
			if artifacts, err = s.ArtifactList(appRef); err != nil {
				return err
			}
			if deployments, err = s.DeploymentList(appRef); err != nil {
				return err
			}
			if releases, err = s.ReleaseList(appRef); err != nil {
				return err
			}
			if statusReports, err = s.StatusReportList(appRef); err != nil {
				return err
			}
		}

		if workspaceList, err := s.ProjectListWorkspaces(ref); err != nil {
			return err
		} else {
			for _, workspace := range workspaceList {
				// Get the triggers for a project in the workspace
				if triggerList, err := s.TriggerList(workspace.Workspace, &pb.Ref_Project{Project: project.Name}, nil, []string{}); err != nil {
					return err
				} else {
					triggers = append(triggers, triggerList...)
				}
				if workspaceDetail, err := s.WorkspaceGet(workspace.Workspace.Workspace); err != nil {
					return err
				} else {
					// If the project we're deleting is the only project in the workspace, we delete the workspace
					// We don't delete the default workspace
					if len(workspaceDetail.Projects) == 1 && workspace.Workspace.Workspace != "default" {
						workspaces = append(workspaces, workspace)
					}
				}

			}
		}

		if pipelines, err = s.PipelineList(ref); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// Builds, artifacts, deployments, releases, status reports, pipelines, triggers,
	// workspaces and config are deleted with a project
	// Jobs and tasks will NOT be deleted along with a project
	// Instances are expected to be deleted before ProjectDelete, via the destroy op
	// delete builds, artifacts, deployments, releases and status reports for each app in the project
	for _, build := range builds {
		if err = s.BuildDelete(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: build.Id}}); err != nil {
			return err
		}
	}
	for _, artifact := range artifacts {
		if err = s.ArtifactDelete(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: artifact.Id}}); err != nil {
			return err
		}
	}
	for _, deployment := range deployments {
		if err = s.DeploymentDelete(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: deployment.Id}}); err != nil {
			return err
		}
	}
	for _, release := range releases {
		if err = s.ReleaseDelete(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: release.Id}}); err != nil {
			return err
		}
	}
	for _, statusReport := range statusReports {
		if err = s.StatusReportDelete(&pb.Ref_Operation{Target: &pb.Ref_Operation_Id{Id: statusReport.Id}}); err != nil {
			return err
		}
	}

	// delete workspaces for a project
	for _, workspace := range workspaces {
		if err = s.WorkspaceDelete(workspace.Workspace.Workspace); err != nil {
			return err
		}
	}

	// delete triggers for project
	for _, trigger := range triggers {
		if err = s.TriggerDelete(&pb.Ref_Trigger{Id: trigger.Id}); err != nil {
			return err
		}
	}

	// delete pipelines for project
	for _, pipeline := range pipelines {
		if err = s.PipelineDelete(&pb.Ref_Pipeline{Ref: &pb.Ref_Pipeline_Id{
			Id: &pb.Ref_PipelineId{Id: pipeline.Id},
		},
		}); err != nil {
			return err
		}
	}

	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	// TODO: Delete config
	err = s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.projectDelete(dbTxn, memTxn, ref)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// ProjectList returns the list of projects.
func (s *State) ProjectList() ([]*pb.Ref_Project, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	return s.projectList(memTxn)
}

// ProjectListWorkspaces returns the list of workspaces that a project is in.
func (s *State) ProjectListWorkspaces(ref *pb.Ref_Project) ([]*pb.Workspace_Project, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	var results []*pb.Workspace_Project
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		results, err = s.workspaceListProjects(dbTxn, memTxn, ref)
		return err
	})

	return results, err
}

// ProjectPollPeek returns the next project that should be polled.
// This will return (nil,nil) if there are no projects to poll currently.
//
// This is a "peek" operation so it does not update the project's next poll
// time. Therefore, calling this multiple times should return the same result
// unless a function like ProjectPollComplete is called.
//
// If ws is non-nil, the WatchSet can be watched for any changes to
// projects to poll. This can be watched, for example, to detect when
// projects to poll are added. This is important functionality since callers
// may be sleeping on a deadline for awhile when a new project is inserted
// to poll immediately.
func (s *State) ProjectPollPeek(ws memdb.WatchSet) (*pb.Project, time.Time, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	// LowerBound doesn't support watches so we have to do a Get first
	// to get a valid watch channel on these fields.
	iter, err := memTxn.Get(
		projectIndexTableName,
		projectIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}
	ws.Add(iter.WatchCh())

	// Get the project with the lowest "next poll" time.
	iter, err = memTxn.LowerBound(
		projectIndexTableName,
		projectIndexNextPollIndexName,
		true,            // polling enabled
		time.Unix(0, 0), // lowest next poll time
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// If we have no values, then return
	raw := iter.Next()
	if raw == nil {
		return nil, time.Time{}, nil
	}

	rec := raw.(*projectIndexRecord)
	if rec.NextPoll.IsZero() {
		// This _shouldnt_ happen but let's protect against it anyways.
		return nil, time.Time{}, nil
	}

	var result *pb.Project
	err = s.db.View(func(dbTxn *bolt.Tx) error {
		var err error
		result, err = s.projectGet(dbTxn, memTxn, &pb.Ref_Project{
			Project: rec.Id,
		})

		return err
	})

	return result, rec.NextPoll, err
}

// ProjectPollComplete sets the next poll time for the given project to the
// time "t" plus the interval time for the project.
func (s *State) ProjectPollComplete(p *pb.Project, t time.Time) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	raw, err := memTxn.First(
		projectIndexTableName,
		projectIndexIdIndexName,
		string(s.projectId(p)),
	)
	if err != nil {
		return err
	}
	if raw == nil {
		return nil
	}

	record := raw.(*projectIndexRecord)
	if !record.Poll {
		// If this project doesn't have polling enabled, then do nothing.
		// This could happen if a project had polling when Peek was called,
		// then between Peek and Complete, polling was disabled.
		return nil
	}

	record = record.Copy()
	record.LastPoll = t
	record.NextPoll = t.Add(record.PollInterval)

	if err := memTxn.Insert(projectIndexTableName, record); err != nil {
		return err
	}

	memTxn.Commit()
	return nil
}

// ProjectUpdateDataRef updates the latest data ref used for a project.
// This data is available via the APIs for querying workspaces.
func (s *State) ProjectUpdateDataRef(
	ref *pb.Ref_Project,
	ws *pb.Ref_Workspace,
	dataRef *pb.Job_DataSource_Ref,
) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.workspaceUpdateProjectDataRef(dbTxn, memTxn, ws, ref, dataRef)
	})
	if err != nil {
		return err
	}

	memTxn.Commit()
	return nil
}

func (s *State) projectGetOrCreate(dbTxn *bolt.Tx, memTxn *memdb.Txn, ref *pb.Ref_Project) (*pb.Project, error) {
	result, err := s.projectGet(dbTxn, memTxn, ref)
	if status.Code(err) == codes.NotFound {
		result = nil
		err = nil
	}
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result = &pb.Project{
		Name: ref.Project,
	}

	return result, s.projectPut(dbTxn, memTxn, result)
}

func (s *State) projectPut(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	value *pb.Project,
) error {
	// This is to prevent mistakes or abuse. Realistically a waypoint.hcl
	// file should be MUCH smaller than this so this catches the really big
	// mistakes.
	if len(value.WaypointHcl) > projectWaypointHclMaxSize {
		return status.Errorf(codes.FailedPrecondition,
			"project 'waypoint_hcl' exceeds maximum size (5MB)",
		)
	}

	id := s.projectId(value)

	// Get the global bucket and write the value to it.
	b := dbTxn.Bucket(projectBucket)
	if err := dbPut(b, id, value); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.projectIndexSet(memTxn, id, value)
}

func (s *State) projectGet(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) (*pb.Project, error) {
	var result pb.Project
	b := dbTxn.Bucket(projectBucket)
	return &result, dbGet(b, []byte(strings.ToLower(ref.Project)), &result)
}

func (s *State) projectList(
	memTxn *memdb.Txn,
) ([]*pb.Ref_Project, error) {
	iter, err := memTxn.Get(projectIndexTableName, projectIndexIdIndexName+"_prefix", "")
	if err != nil {
		return nil, err
	}

	var result []*pb.Ref_Project
	for {
		next := iter.Next()
		if next == nil {
			break
		}
		idx := next.(*projectIndexRecord)

		result = append(result, &pb.Ref_Project{
			Project: idx.Id,
		})
	}

	return result, nil
}

func (s *State) projectDelete(
	dbTxn *bolt.Tx,
	memTxn *memdb.Txn,
	ref *pb.Ref_Project,
) error {
	// Get the project. If it doesn't exist then we're successful.
	p, err := s.projectGet(dbTxn, memTxn, ref)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil
		}

		return err
	}

	// Delete each application
	for _, app := range p.Applications {
		if err := s.appDelete(dbTxn, memTxn, &pb.Ref_Application{
			Project:     ref.Project,
			Application: app.Name,
		}); err != nil {
			return err
		}
	}

	// Delete from bolt
	id := s.projectIdByRef(ref)
	if err := dbTxn.Bucket(projectBucket).Delete(id); err != nil {
		return err
	}

	// Delete from memdb
	if err := memTxn.Delete(projectIndexTableName, &projectIndexRecord{Id: string(id)}); err != nil {
		return err
	}

	return nil
}

// projectIndexSet writes an index record for a single project.
func (s *State) projectIndexSet(txn *memdb.Txn, id []byte, value *pb.Project) error {
	record := &projectIndexRecord{
		Id:            string(id),
		Poll:          false, // being explicit that we want to default poll to false
		AppStatusPoll: false, // application polling off by default until a deployment or release happens
	}

	// This entire if block sets up polling tracking for the project. In the
	// state store we just maintain timestamps of when to poll next. It is
	// up to downstream users to call ProjectNextPoll repeatedly to iterate
	// over the next projects to poll and do something.
	if p := value.DataSourcePoll; p != nil && p.Enabled {
		// If it's empty at this point, we'll set the default here.
		if p.Interval == "" {
			p.Interval = defaultProjectPollInterval
		}
		interval, err := time.ParseDuration(p.Interval)
		if err != nil {
			return err
		}

		// We're polling. By default we have no last polling time and
		// we set the next polling time to now cause we want to poll ASAP.
		// If we're updating a project without changing the poll settings,
		// the next block will ensure we have the next poll time retained.
		record.Poll = true
		record.NextPoll = time.Now()
		record.PollInterval = interval

		// If there is a previous value with a last poll time, then we
		// update the next poll time to use our new interval.
		raw, err := txn.First(
			projectIndexTableName,
			projectIndexIdIndexName,
			record.Id,
		)
		if err != nil {
			return err
		}
		if raw != nil {
			recordOld := raw.(*projectIndexRecord)

			// If we have a last poll time, then set the next poll time.
			// This also ensures that if we're updating a project w/o changing
			// poll settings, that the previous settings are retained.
			if !recordOld.LastPoll.IsZero() {
				record.LastPoll = recordOld.LastPoll
				record.NextPoll = record.LastPoll.Add(interval)
			}
		}
	}

	// Insert application poll

	// Note that application status polling currently only turned on by default
	// if project polling is enabled. This is because application status polling
	// requires a data source, which currently is only possible through git polling.
	app := value.StatusReportPoll
	if app == nil && record.Poll {
		// Auto-turn on app polling. Can be disabled if project settings explicitly
		// disable app status polling
		value.StatusReportPoll = &pb.Project_AppStatusPoll{
			Enabled: true,
		}
	}

	// This entire if block sets up polling tracking for the application. In the
	// state store we just maintain timestamps of when to poll next. It is
	// up to downstream users to call ApplicationNextPoll repeatedly to iterate
	// over the next projects to poll and do something.
	if app := value.StatusReportPoll; app != nil && app.Enabled {
		if app.Interval == "" {
			app.Interval = defaultAppStatusPollInterval
		}
		interval, err := time.ParseDuration(app.Interval)
		if err != nil {
			return err
		}

		// We're polling. By default we have no last polling time and
		// we set the next polling time to now cause we want to poll ASAP.
		// If we're updating an app without changing the poll settings,
		// the next block will ensure we have the next poll time retained.
		record.AppStatusPoll = true
		record.AppStatusNextPoll = time.Now()
		record.AppStatusPollInterval = interval

		// If there is a previous value with a last poll time, then we
		// update the next poll time to use our new interval.
		raw, err := txn.First(
			projectIndexTableName,
			projectIndexIdIndexName,
			record.Id,
		)
		if err != nil {
			return err
		}
		if raw != nil {
			recordOld := raw.(*projectIndexRecord)

			// If we have a last poll time, then set the next poll time.
			// This also ensures that if we're updating an app w/o changing
			// poll settings, that the previous settings are retained.
			if !recordOld.AppStatusLastPoll.IsZero() {
				record.AppStatusLastPoll = recordOld.AppStatusLastPoll
				record.AppStatusNextPoll = record.AppStatusLastPoll.Add(interval)
			}
		}
	}

	// Insert the index
	return txn.Insert(projectIndexTableName, record)
}

// projectIndexInit initializes the project index from persisted data.
func (s *State) projectIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(projectBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var value pb.Project
		if err := proto.Unmarshal(v, &value); err != nil {
			return err
		}
		if err := s.projectIndexSet(memTxn, k, &value); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) projectId(p *pb.Project) []byte {
	return []byte(strings.ToLower(p.Name))
}

func (s *State) projectIdByRef(ref *pb.Ref_Project) []byte {
	return []byte(strings.ToLower(ref.Project))
}

func projectIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: projectIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			projectIndexIdIndexName: {
				Name:         projectIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field:     "Id",
					Lowercase: true,
				},
			},

			projectIndexNextPollIndexName: {
				Name:         projectIndexNextPollIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "Poll",
						},

						&IndexTime{
							Field: "NextPoll",
							Asc:   true,
						},
					},
				},
			},

			appIndexNextPollIndexName: {
				Name:         appIndexNextPollIndexName,
				AllowMissing: true,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.BoolFieldIndex{
							Field: "AppStatusPoll",
						},

						&IndexTime{
							Field: "AppStatusNextPoll",
							Asc:   true,
						},
					},
				},
			},
		},
	}
}

const (
	projectIndexTableName         = "project-index"
	projectIndexIdIndexName       = "id"
	projectIndexNextPollIndexName = "next-poll"
	appIndexNextPollIndexName     = "app-next-poll"

	projectWaypointHclMaxSize = 5 * 1024 // 5 MB

	// defaultProjectPollInterval is used by the project poll handler
	// for setting up a default interval time
	defaultProjectPollInterval = "30s"

	// defaultAppStatusPollInterval is used for polling a status report job for each
	// application defined in a project. It is initially set to a long interval
	// so that Waypoint doesn't overrun and rate limit user accounts like on AWS.
	// Users must opt into shorter interval times.
	defaultAppStatusPollInterval = "5m"
)

type projectIndexRecord struct {
	Id string

	// Project polling is used for updating the project from a remote source
	// on an interval

	// Poll is true if this project has polling enabled.
	Poll bool
	// PollInterval is the interval currently set between poll operations.
	PollInterval time.Duration
	// LastPoll is the time that the last polling operation was queued.
	// NextPoll is the time when the next polling operation is expected.
	// Storing NextPoll rather than the interval makes it easier to query
	// for the next project.
	LastPoll time.Time
	NextPoll time.Time

	// Application Polling is used for generating a status report on the current
	// health of all applications in a project.

	// AppStatusPoll is true if this projects applications has polling enabled.
	AppStatusPoll bool
	// AppStatusPollInterval is the interval currently set between poll operations.
	AppStatusPollInterval time.Duration
	// We separate project and application polling vars because project polling
	// is used for updating the project, and application polling is used for
	// generating status reports. So there are two separate Next and Last Poll
	// vars for projects and applications
	AppStatusLastPoll time.Time
	AppStatusNextPoll time.Time
}

// Copy should be called prior to any modifications to an existing record.
func (idx *projectIndexRecord) Copy() *projectIndexRecord {
	// A shallow copy is good enough since we only modify top-level fields.
	copy := *idx
	return &copy
}
