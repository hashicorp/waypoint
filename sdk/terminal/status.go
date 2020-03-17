package terminal

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

// Status is used to provide an updating status to the user. The status
// usually has some animated element along with it such as a spinner.
type Status interface {
	// Update writes a new status. This should be a single line.
	Update(string)

	// Close should be called when the live updating is complete. The
	// status will be cleared from the line.
	Close() error
}

// spinnerStatus implements Status and uses a spinner to show updates.
type spinnerStatus struct {
	mu      sync.Mutex
	spinner *spinner.Spinner
}

func newSpinnerStatus() *spinnerStatus {
	return &spinnerStatus{
		spinner: spinner.New(
			spinner.CharSets[11],
			time.Second/6,
			spinner.WithColor("bold"),
		),
	}
}

func (s *spinnerStatus) Update(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.spinner.Suffix = " " + msg
}

func (s *spinnerStatus) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.spinner.Stop()
	return nil
}
