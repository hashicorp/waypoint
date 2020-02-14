package status

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

type SpinnerStatus struct {
	mu      sync.Mutex
	once    sync.Once
	spinner *spinner.Spinner
}

func (s *SpinnerStatus) Update(str string) {
	s.once.Do(func() {
		if s.spinner == nil {
			s.spinner = spinner.New(spinner.CharSets[11], time.Second/6, spinner.WithSuffix(" "+str), spinner.WithColor("bold"))
		}
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	// To be sure we don't crash if someone tries to status update after we've closed
	if s.spinner != nil {
		s.spinner.Suffix = " " + str
	}
}

func (s *SpinnerStatus) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.spinner != nil {
		s.spinner.Stop()
	}
}
