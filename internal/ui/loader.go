package ui

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Loader displays a simple spinner while work is in progress.
type Loader struct {
	out    io.Writer
	ticker *time.Ticker
	done   chan struct{}
	frames []string
}

// NewLoader creates a new Loader with default spinner frames.
func NewLoader() *Loader {
	return &Loader{
		out:    os.Stderr,
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start begins the spinner.
func (l *Loader) Start() {
	if l.ticker != nil {
		return
	}
	l.ticker = time.NewTicker(80 * time.Millisecond)
	t := l.ticker
	l.done = make(chan struct{})
	go func() {
		i := 0
		for {
			select {
			case <-t.C:
				fmt.Fprintf(l.out, "\r%s", l.frames[i%len(l.frames)])
				i++
			case <-l.done:
				fmt.Fprint(l.out, "\r")
				return
			}
		}
	}()
}

// Stop stops the spinner and moves to the next line.
func (l *Loader) Stop() {
	if l.ticker == nil {
		return
	}
	l.ticker.Stop()
	close(l.done)
	fmt.Fprintln(l.out)
}
