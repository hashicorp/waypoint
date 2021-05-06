package logbuffer

import (
	"io"
	"time"
)

// LogMerge can combine multiple log streams into one stream. It presumes
// each stream emits entries in time order, and then weaves the entries
// together to create a total time ordered stream.
type LogMerge struct {
	inputs []LogMergeInput
	heads  []TimedEntry
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

// LogMergeInput is value that returns TimedEntry's for LogMerge
// to weave together.
type LogMergeInput interface {
	Next() (TimedEntry, error)
}

// NewLogMerge creates a new LogMerge, with the stream generated
// from the given inputs.
func NewLogMerge(inputs ...LogMergeInput) (*LogMerge, error) {
	var lm LogMerge
	lm.inputs = inputs

	lm.heads = lm.makeHeads()

	return &lm, nil
}

// AddInput adds a new input to the LogMerge. NOTE: inputs can
// not be removed at this time.
func (l *LogMerge) AddInput(input LogMergeInput) {
	l.inputs = append(l.inputs, input)
	l.heads = append(l.heads, nil)
}

// TimedEntries is a convience type of TimedEntry's that provides
// the LogMergeInput interface.
type TimedEntries []TimedEntry

// Next returns the next value in the slice and then shrinks itself.
func (t *TimedEntries) Next() (TimedEntry, error) {
	if len(*t) == 0 {
		return nil, io.EOF
	}

	ent := (*t)[0]

	*t = (*t)[1:]

	return ent, nil
}

// Create a slice to be used by refillEntries and findNext
func (l *LogMerge) makeHeads() []TimedEntry {
	return make([]TimedEntry, len(l.inputs))
}

// Populate any missing entries with an entry from the corresponding
// input (index in slice corresponds to inputs slice).
func (l *LogMerge) refillEntries(entries []TimedEntry) int {
	var pop int

	for i, ent := range entries {
		if ent != nil {
			pop++
			continue
		}

		ent, err := l.inputs[i].Next()
		if err == nil {
			pop++
			entries[i] = ent
		}
	}

	return pop
}

// Find the entry in entries with earliest time, returning the entry and the
// input that generated it.
func (l *LogMerge) findNext(entries []TimedEntry) (TimedEntry, LogMergeInput) {
	var (
		best      TimedEntry
		bestIdx   int
		bestInput LogMergeInput
	)

	for i, ent := range entries {
		if ent == nil {
			continue
		}

		if best == nil || ent.Time().Before(best.Time()) {
			best = ent
			bestIdx = i
			bestInput = l.inputs[i]
		}
	}

	if best == nil {
		return nil, nil
	}

	entries[bestIdx] = nil

	return best, bestInput
}

// InputEntry is returned by ReadNext. It provides access to the TimedEntry
// that is next as well as the input that generated the entry. This
// type is important because it allows the caller to figure out the context
// of the entry from the input. Because LogMerge is going to effectively
// shuffle the values that are put into it, the caller is going to have to deal
// with entries appearing in any order and the input provides critical context.
type InputEntry struct {
	TimedEntry
	Input LogMergeInput
}

// ReadNext returns a slice of InputEntrys that are next in total time order.
// The result might be fewer than count values, depending on what is available.
func (l *LogMerge) ReadNext(count int) ([]InputEntry, error) {
	var out []InputEntry

	heads := l.heads

	for i := 0; i < count; i++ {
		pop := l.refillEntries(heads)
		if pop == 0 {
			break
		}

		ent, ip := l.findNext(heads)
		out = append(out, InputEntry{TimedEntry: ent, Input: ip})
	}

	return out, nil
}
