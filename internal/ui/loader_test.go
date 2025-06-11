package ui

import (
	"bytes"
	"testing"
	"time"
)

func TestLoaderStartStop(t *testing.T) {
	l := NewLoader()
	var buf bytes.Buffer
	l.s.Writer = &buf
	l.Start()
	time.Sleep(90 * time.Millisecond)
	l.Stop()
	if buf.Len() == 0 {
		t.Errorf("expected spinner output")
	}
}
