package state

import (
	"math"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

var buildBucket = []byte("build")

const (
	buildIndexTableName             = "build-index"
	buildIndexIdIndexName           = "id"
	buildIndexStartTimeIndexName    = "start-time-by-app"
	buildIndexCompleteTimeIndexName = "complete-time-by-app"
)

func init() {
	dbBuckets = append(dbBuckets, buildBucket)
	dbIndexers = append(dbIndexers, (*State).buildIndexInit)
	schemas = append(schemas, buildIndexSchema)
}

func buildIndexSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: buildIndexTableName,
		Indexes: map[string]*memdb.IndexSchema{
			buildIndexIdIndexName: &memdb.IndexSchema{
				Name:         buildIndexIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Id",
				},
			},

			buildIndexStartTimeIndexName: &memdb.IndexSchema{
				Name:         buildIndexStartTimeIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&IndexTime{
							Field: "StartTime",
						},
					},
				},
			},

			buildIndexCompleteTimeIndexName: &memdb.IndexSchema{
				Name:         buildIndexCompleteTimeIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field:     "Project",
							Lowercase: true,
						},

						&memdb.StringFieldIndex{
							Field:     "App",
							Lowercase: true,
						},

						&IndexTime{
							Field: "CompleteTime",
						},
					},
				},
			},
		},
	}
}

type buildIndexRecord struct {
	Id           string
	Project      string
	App          string
	StartTime    time.Time
	CompleteTime time.Time
}

// MatchRef checks if a record matches the ref value. We have to provide
// this because we use LowerBound lookups in memdb and this may return
// a non-matching value at a certain point.
func (rec *buildIndexRecord) MatchRef(ref *pb.Ref_Application) bool {
	return rec.Project == ref.Project && rec.App == ref.Application
}

// BuildPut inserts or updates a build record.
func (s *State) BuildPut(update bool, b *pb.Build) error {
	memTxn := s.inmem.Txn(true)
	defer memTxn.Abort()

	err := s.db.Update(func(dbTxn *bolt.Tx) error {
		return s.buildPut(dbTxn, memTxn, update, b)
	})
	if err == nil {
		memTxn.Commit()
	}

	return err
}

// BuildGet gets a build by ID.
func (s *State) BuildGet(id string) (*pb.Build, error) {
	var result pb.Build
	err := s.db.View(func(tx *bolt.Tx) error {
		return dbGet(tx.Bucket(buildBucket), []byte(id), &result)
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *State) BuildList(ref *pb.Ref_Application) ([]*pb.Build, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.LowerBound(
		buildIndexTableName,
		buildIndexStartTimeIndexName,
		ref.Project,
		ref.Application,
		time.Unix(math.MaxInt64, 0),
	)
	if err != nil {
		return nil, err
	}

	var result []*pb.Build
	s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(buildBucket)

		for {
			current := iter.Next()
			if current == nil {
				return nil
			}

			record := current.(*buildIndexRecord)
			if !record.MatchRef(ref) {
				return nil
			}

			var build pb.Build
			if err := dbGet(bucket, []byte(record.Id), &build); err != nil {
				return err
			}

			result = append(result, &build)
		}
	})

	return result, nil
}

// BuildLatest gets the latest build that was completed successfully.
func (s *State) BuildLatest(ref *pb.Ref_Application) (*pb.Build, error) {
	memTxn := s.inmem.Txn(false)
	defer memTxn.Abort()

	iter, err := memTxn.LowerBound(
		buildIndexTableName,
		buildIndexCompleteTimeIndexName,
		ref.Project,
		ref.Application,
		time.Unix(math.MaxInt64, 0),
	)
	if err != nil {
		return nil, err
	}

	for {
		raw := iter.Next()
		if raw == nil {
			return nil, nil
		}

		record := raw.(*buildIndexRecord)
		if !record.MatchRef(ref) {
			return nil, nil
		}

		b, err := s.BuildGet(record.Id)
		if err != nil {
			return nil, err
		}

		// Shouldn't happen but if it does, return nothing.
		if b.Status == nil {
			return nil, nil
		}

		// State must be success.
		switch b.Status.State {
		case pb.Status_SUCCESS:
			return b, nil
		}
	}
}

func (s *State) buildPut(
	tx *bolt.Tx,
	inmemTxn *memdb.Txn,
	update bool,
	build *pb.Build,
) error {
	id := []byte(build.Id)

	// Get the global bucket and write the value to it.
	b := tx.Bucket(buildBucket)
	if err := dbPut(b, id, build); err != nil {
		return err
	}

	// Create our index value and write that.
	return s.buildPutIndex(inmemTxn, build)
}

func (s *State) buildIndexInit(dbTxn *bolt.Tx, memTxn *memdb.Txn) error {
	bucket := dbTxn.Bucket(buildBucket)
	return bucket.ForEach(func(k, v []byte) error {
		var build pb.Build
		if err := proto.Unmarshal(v, &build); err != nil {
			return err
		}
		if err := s.buildPutIndex(memTxn, &build); err != nil {
			return err
		}

		return nil
	})
}

func (s *State) buildPutIndex(txn *memdb.Txn, build *pb.Build) error {
	var startTime, completeTime time.Time
	if build.Status != nil {
		if t := build.Status.StartTime; t != nil {
			st, err := ptypes.Timestamp(t)
			if err != nil {
				return status.Errorf(codes.Internal, "time for build can't be parsed")
			}

			startTime = st
		}

		if t := build.Status.CompleteTime; t != nil {
			ct, err := ptypes.Timestamp(build.Status.CompleteTime)
			if err != nil {
				return status.Errorf(codes.Internal, "time for build can't be parsed")
			}

			completeTime = ct
		}
	}

	return txn.Insert(buildIndexTableName, &buildIndexRecord{
		Id:           build.Id,
		Project:      build.Application.Project,
		App:          build.Application.Application,
		StartTime:    startTime,
		CompleteTime: completeTime,
	})
}
