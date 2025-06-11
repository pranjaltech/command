package ui

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Loader displays a simple spinner while work is in progress.
type Loader struct {
	msg    string
	out    io.Writer
	ticker *time.Ticker
	done   chan struct{}
	frames []rune
}

// NewLoader creates a new Loader with default spinner frames.
func NewLoader(msg string) *Loader {
	return &Loader{
		msg:    msg,
		out:    os.Stderr,
		frames: []rune{'|', '/', '-', '\\'},
	}
}

// Start begins the spinner.
func (l *Loader) Start() {
	if l.ticker != nil {
		return
	}
	l.ticker = time.NewTicker(120 * time.Millisecond)
	t := l.ticker
	l.done = make(chan struct{})
	go func() {
		i := 0
		for {
			select {
			case <-t.C:
				fmt.Fprintf(l.out, "\r%s %c", l.msg, l.frames[i%len(l.frames)])
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
