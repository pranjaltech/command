package ui

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

// Loader wraps a spinner that shows progress during long operations.
type Loader struct{ s *spinner.Spinner }

// NewLoader returns a loader using a random character set and 80ms delay.
func NewLoader() *Loader {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	keys := make([]int, 0, len(spinner.CharSets))
	for k := range spinner.CharSets {
		keys = append(keys, k)
	}
	cs := spinner.CharSets[keys[r.Intn(len(keys))]]
	sp := spinner.New(cs, 80*time.Millisecond, spinner.WithWriter(os.Stderr))
	return &Loader{s: sp}
}

// Start begins the spinner.
func (l *Loader) Start() { l.s.Start() }

// Stop stops the spinner and moves to the next line.
func (l *Loader) Stop() {
	l.s.Stop()
	fmt.Fprint(l.s.Writer, "\r")
}
