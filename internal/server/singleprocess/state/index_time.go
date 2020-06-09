package state

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"
)

// IndexTime indexes a time.Time field of a struct.
// TODO(mitchellh): test
type IndexTime struct {
	Field string

	// Asc if true will index with ascending order. Usually you want to
	// find the most recent data so this is false by default.
	Asc bool
}

func (idx *IndexTime) FromObject(obj interface{}) (bool, []byte, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	fv := v.FieldByName(idx.Field)
	if !fv.IsValid() {
		return false, nil,
			fmt.Errorf("field '%s' is invalid %#v ", idx.Field, obj)
	}

	timeVal, ok := fv.Interface().(time.Time)
	if !ok {
		return false, nil,
			fmt.Errorf("field '%s' is not a time %v ", idx.Field, obj)
	}

	return true, idx.fromTime(timeVal), nil
}

func (idx *IndexTime) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single argument")
	}

	arg, ok := args[0].(time.Time)
	if !ok {
		return nil, fmt.Errorf("argument must be a time: %#v", args[0])
	}

	return idx.fromTime(arg), nil
}

func (idx *IndexTime) fromTime(t time.Time) []byte {
	// If we're ascending, we just use the unix timestamp. This will give us
	// the proper ordering.
	val := t.UnixNano()

	// If we're descending, we use the time until max int64 which is the
	// maximum size of a Go unix timestamp. This will give us a SMALLER
	// value for NEWER (more recent) times. A descending order.
	if !idx.Asc {
		val = math.MaxInt64 - val
	}

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, val)

	// The endianness of binary.PutVarint is reversed for byte-based
	// sorting. We need to reverse.
	for i := len(buf)/2 - 1; i >= 0; i-- {
		opp := len(buf) - 1 - i
		buf[i], buf[opp] = buf[opp], buf[i]
	}

	return buf
}
