// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logbuffer

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testTimeEntry struct {
	msg string
	ts  time.Time
}

func (t *testTimeEntry) Time() time.Time {
	return t.ts
}

func (t *testTimeEntry) Value() interface{} {
	return t
}

func TestMerge(t *testing.T) {
	t.Run("returns entries directly if only one source", func(t *testing.T) {
		var input TimedEntries

		start := time.Now()

		for i := 0; i < 10; i++ {
			input = append(input, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})
		}

		lm := NewMerger(&input)

		entries, err := lm.Read(10)
		require.NoError(t, err)

		require.Len(t, entries, 10)

		for i := 0; i < 10; i++ {
			require.True(t, start.Add(time.Duration(i)*time.Second).Equal(entries[i].Time()))
		}
	})

	t.Run("interleaves entries from multiple inputs", func(t *testing.T) {
		var (
			input1 TimedEntries
			input2 TimedEntries
		)

		start := time.Now()

		for i := 0; i < 10; i++ {
			input1 = append(input1, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})

			input2 = append(input2, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})
		}

		lm := NewMerger(&input1, &input2)

		entries, err := lm.Read(20)
		require.NoError(t, err)

		require.Len(t, entries, 20)

		var sec time.Duration
		for i := 0; i < 20; i += 2 {
			require.True(t, start.Add(sec).Equal(entries[i].Time()))
			require.True(t, start.Add(sec).Equal(entries[i+1].Time()))
			sec += time.Second
		}
	})

	t.Run("doesn't lose entries if there isn't enough space", func(t *testing.T) {
		var (
			input1 TimedEntries
			input2 TimedEntries
		)

		start := time.Now()

		for i := 0; i < 10; i++ {
			input1 = append(input1, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})

			input2 = append(input2, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})
		}

		lm := NewMerger(&input1, &input2)

		entries, err := lm.Read(1)
		require.NoError(t, err)

		require.Len(t, entries, 1)

		require.True(t, start.Equal(entries[0].Time()))

		entries, err = lm.Read(1)
		require.NoError(t, err)

		require.Len(t, entries, 1)
		require.True(t, start.Equal(entries[0].Time()))
	})

	t.Run("deals with inputs have different total entries", func(t *testing.T) {
		var (
			input1 TimedEntries
			input2 TimedEntries
		)

		start := time.Now()

		for i := 0; i < 10; i++ {
			input1 = append(input1, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})

			if i < 5 {
				input2 = append(input2, &testTimeEntry{
					msg: fmt.Sprintf("entry %d", i),
					ts:  start.Add(time.Duration(i) * time.Second),
				})
			}
		}

		lm := NewMerger(&input1, &input2)

		entries, err := lm.Read(20)
		require.NoError(t, err)

		require.Len(t, entries, 15)

		var sec time.Duration
		for i := 0; i < 10; i += 2 {
			require.True(t, start.Add(sec).Equal(entries[i].Time()))
			if i < 5 {
				require.True(t, start.Add(sec).Equal(entries[i+1].Time()))
			}
			sec += time.Second
		}

		for i := 10; i < 15; i++ {
			require.True(t, start.Add(sec).Equal(entries[i].Time()))
			sec += time.Second
		}
	})

	t.Run("handles disjoint inputs", func(t *testing.T) {
		var (
			input1 TimedEntries
			input2 TimedEntries
		)

		start := time.Now()

		for i := 0; i < 10; i++ {
			input1 = append(input1, &testTimeEntry{
				msg: fmt.Sprintf("entry %d", i),
				ts:  start.Add(time.Duration(i) * time.Second),
			})
		}

		for i := 0; i < 10; i++ {
			input2 = append(input2, &testTimeEntry{
				msg: fmt.Sprintf("entry 2 %d", i),
				ts:  start.Add(time.Duration(100*(i+1)) * time.Second),
			})
		}

		lm := NewMerger(&input1, &input2)

		entries, err := lm.Read(20)
		require.NoError(t, err)

		require.Len(t, entries, 20)

		for i := 0; i < 10; i++ {
			require.True(t, start.Add(time.Duration(i)*time.Second).Equal(entries[i].Time()))
		}

		for i := 0; i < 10; i++ {
			require.True(t, start.Add(time.Duration(100*(i+1))*time.Second).Equal(entries[i+10].Time()))
		}
	})

}
