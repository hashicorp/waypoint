// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	"sync/atomic"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/pkg/server/logbuffer"
)

const (
	instanceLogsTableName           = "instance-logs"
	instanceLogsIdIndexName         = "id"
	instanceLogsInstanceIdIndexName = "deployment-id"
)

func init() {
	schemas = append(schemas, instanceLogsSchema)
}

var instanceLogsId int64

func instanceLogsSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: instanceLogsTableName,
		Indexes: map[string]*memdb.IndexSchema{
			instanceLogsIdIndexName: {
				Name:         instanceLogsIdIndexName,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.IntFieldIndex{
					Field: "Id",
				},
			},

			instanceLogsInstanceIdIndexName: {
				Name:         instanceLogsInstanceIdIndexName,
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

// InstanceLogs is a value that can be created to assist in coordination
// log writers and readers. It is a lighter weight version of an Instance
// used to manage virtual CEBs sending logs.
type InstanceLogs struct {
	Id         int64
	InstanceId string

	LogBuffer *logbuffer.Buffer
}

func (s *State) InstanceLogsCreate(ctx context.Context, id string, logs *InstanceLogs) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()

	// Set our ID
	logs.Id = atomic.AddInt64(&instanceLogsId, 1)
	logs.InstanceId = id

	// Insert
	if err := txn.Insert(instanceLogsTableName, logs); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceLogsDelete(ctx context.Context, id int64) error {
	txn := s.inmem.Txn(true)
	defer txn.Abort()
	if _, err := txn.DeleteAll(instanceLogsTableName, instanceLogsIdIndexName, id); err != nil {
		return status.Errorf(codes.Aborted, err.Error())
	}
	txn.Commit()

	return nil
}

func (s *State) InstanceLogsById(ctx context.Context, id int64) (*InstanceLogs, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceLogsTableName, instanceLogsIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance exec ID not found")
	}

	return raw.(*InstanceLogs), nil
}

func (s *State) InstanceLogsByInstanceId(ctx context.Context, id string) (*InstanceLogs, error) {
	txn := s.inmem.Txn(false)
	raw, err := txn.First(instanceLogsTableName, instanceLogsInstanceIdIndexName, id)
	txn.Abort()
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, status.Errorf(codes.NotFound, "instance exec ID not found")
	}

	return raw.(*InstanceLogs), nil
}

func (s *State) InstanceLogsListByInstanceId(ctx context.Context, id string, ws memdb.WatchSet) ([]*InstanceLogs, error) {
	txn := s.inmem.Txn(false)
	defer txn.Abort()
	return s.instanceLogsListByInstanceId(txn, id, ws)
}

func (s *State) instanceLogsListByInstanceId(
	txn *memdb.Txn, id string, ws memdb.WatchSet,
) ([]*InstanceLogs, error) {
	// Find all the exec sessions
	iter, err := txn.Get(instanceLogsTableName, instanceLogsInstanceIdIndexName, id)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, err.Error())
	}

	var result []*InstanceLogs
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		result = append(result, raw.(*InstanceLogs))
	}

	ws.Add(iter.WatchCh())

	return result, nil
}
