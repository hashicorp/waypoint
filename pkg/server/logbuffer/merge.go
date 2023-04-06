// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logbuffer

import (
	"io"
	"time"
)

// Merger can combine multiple log streams into one stream. It presumes
// each stream emits entries in time order, and then weaves the entries
// together to create a total time ordered stream.
type Merger struct {
	readers []MergeReader
	heads   []TimedEntry
}

// TimedEntry is the interface each input returns entries in
type TimedEntry interface {
	// Time is the associated time of the value. This is used to sort
	// the entry against other entries from other inputs.
	Time() time.Time

	// Value returns the log entry itself. This allows time addition
	// wrappers to return their original value.
	Value() interface{}
}

// MergeReader is value that returns TimedEntry's for Merger
// to weave together.
type MergeReader interface {
	NextTimedEntry() (TimedEntry, error)
}

// NewMerger creates a new Merger, with the stream generated
// from the given inputs.
func NewMerger(readers ...MergeReader) *Merger {
	return &Merger{
		readers: readers,
		heads:   make([]TimedEntry, len(readers)),
	}
}

// TimedEntries is a convience type of TimedEntry's that provides
// the MergeReader interface.
type TimedEntries []TimedEntry

// Next returns the next value in the slice and then shrinks itself.
func (t *TimedEntries) NextTimedEntry() (TimedEntry, error) {
	if len(*t) == 0 {
		return nil, io.EOF
	}

	ent := (*t)[0]

	*t = (*t)[1:]

	return ent, nil
}

// Populate any missing entries with an entry from the corresponding
// input (index in slice corresponds to inputs slice).
func (l *Merger) refillEntries(entries []TimedEntry) int {
	var pop int

	for i, ent := range entries {
		if ent != nil {
			pop++
			continue
		}

		ent, err := l.readers[i].NextTimedEntry()
		if err == nil {
			pop++
			entries[i] = ent
		}
	}

	return pop
}

// Find the entry in entries with earliest time, returning the entry and the
// input that generated it.
func (l *Merger) findNext(entries []TimedEntry) (TimedEntry, MergeReader) {
	var (
		best      TimedEntry
		bestIdx   int
		bestInput MergeReader
	)

	for i, ent := range entries {
		if ent == nil {
			continue
		}

		if best == nil || ent.Time().Before(best.Time()) {
			best = ent
			bestIdx = i
			bestInput = l.readers[i]
		}
	}

	if best == nil {
		return nil, nil
	}

	entries[bestIdx] = nil

	return best, bestInput
}

// ReaderEntry is returned by ReadNext. It provides access to the TimedEntry
// that is next as well as the input that generated the entry. This
// type is important because it allows the caller to figure out the context
// of the entry from the input. Because Merger is going to effectively
// shuffle the values that are put into it, the caller is going to have to deal
// with entries appearing in any order and the input provides critical context.
type ReaderEntry struct {
	TimedEntry
	Reader MergeReader
}

// ReadNext returns a slice of InputEntrys that are next in total time order.
// The result might be fewer than count values, depending on what is available.
func (l *Merger) Read(count int) ([]ReaderEntry, error) {
	var out []ReaderEntry

	heads := l.heads

	for i := 0; i < count; i++ {
		pop := l.refillEntries(heads)
		if pop == 0 {
			break
		}

		ent, ip := l.findNext(heads)
		out = append(out, ReaderEntry{TimedEntry: ent, Reader: ip})
	}

	return out, nil
}
