package state

import (
	"io"

	"github.com/hashicorp/go-memdb"
	/*
		"google.golang.org/grpc/codes"
		"google.golang.org/grpc/status"
	*/

	pb "github.com/mitchellh/devflow/internal/server/gen"
)

const (
	instanceExecTableName           = "instance-execs"
	instanceExecIdIndexName         = "id"
	instanceExecInstanceIdIndexName = "deployment-id"
)

func init() {
	schemas = append(schemas, instanceExecSchema)
}

func instanceExecSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: instanceExecTableName,
		Indexes: map[string]*memdb.IndexSchema{
			instanceExecIdIndexName: &memdb.IndexSchema{
				Name:         instanceExecIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.IntFieldIndex{
					Field: "Id",
				},
			},

			instanceExecInstanceIdIndexName: &memdb.IndexSchema{
				Name:         instanceExecInstanceIdIndexName,
				AllowMissing: false,
				Unique:       false,
				Indexer: &memdb.StringFieldIndex{
					Field:     "InstanceId",
					Lowercase: true,
				},
			},
		},
	}
}

type InstanceExec struct {
	Id         int64
	InstanceId string
	Args       []string
	Reader     io.Reader
	EventCh    chan<- *pb.EntrypointExecRequest
	Connected  uint32
}
